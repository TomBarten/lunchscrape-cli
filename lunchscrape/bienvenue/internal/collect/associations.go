package collect

import (
    "fmt"
    "strconv"
    "strings"

    "github.com/PuerkitoBio/goquery"
    "github.com/TomBarten/lunchscrape_cli/models"
    "github.com/gocolly/colly"
)

func collectItemAssociations(element *colly.HTMLElement, currencySymbol string) *[]models.ItemAssociation {

    associationGroups := element.DOM.Find(
        "form#product-form div.product-section.product-section-input fieldset.product-association-group")

    if len(associationGroups.Nodes) <= 0 {
        empty := make([]models.ItemAssociation, 0)
        return &empty
    }

    associations := make([]models.ItemAssociation, 0, 10)

    associationGroups.EachWithBreak(func(groupIteration int, groupSelection *goquery.Selection) bool {

        isOptional := false

        optionalSpan := groupSelection.Find("legend span.badge.control-optional")

        if optionalSpan != nil && len(optionalSpan.Nodes) > 0 {
            isOptional = true
        }

        error := handleAssociationGroup(&associations, groupSelection, currencySymbol, isOptional)

        if error == nil {
            return true
        }

        fmt.Println(error)
        return false
    })

    return &associations
}

func handleAssociationGroup(
    associations *[]models.ItemAssociation,
    groupSelection *goquery.Selection,
    currencySymbol string,
    isOptional bool) error {

    associationNameLegends := groupSelection.Find("legend:not([class])")

    if len(associationNameLegends.Nodes) <= 0 {
        return fmt.Errorf("cannot find association name element")
    }

    associationName := strings.TrimSpace(associationNameLegends.First().Contents().Not("span").Text())

    if len(associationName) <= 0 {
        return fmt.Errorf("association name is empty")
    }

    associationElements := groupSelection.Find(
        "div.product-association-group-associations label.product-association")

    associationGroupInputs := groupSelection.ChildrenFiltered(":input[id$=\"__Id\"][type=hidden]")

    if len(associationGroupInputs.Nodes) <= 0 {
        return fmt.Errorf("cannot find association group identifier element")
    }

    groupId := associationGroupInputs.First().AttrOr("value", "")

    if len(groupId) <= 0 {
        return fmt.Errorf("association group identifier is empty")
    }

    associationElements.EachWithBreak(func(i int, associationSelection *goquery.Selection) bool {

        association, error := constructItemAssociation(
            associationName, groupId, isOptional, currencySymbol, associationSelection)

        if error != nil {
            fmt.Println(error)
            return false
        }

        *associations = append(*associations, *association)
        return true
    })

    return nil
}

func constructItemAssociation(
    name string,
    groupId string,
    isOptional bool,
    currencySymbol string,
    associationSelection *goquery.Selection) (*models.ItemAssociation, error) {

    isAlwaysChecked := true

    if associationSelection.HasClass("checkbox") || !associationSelection.HasClass("checked") {
        isAlwaysChecked = false
    }

    associationIdInputs := associationSelection.ChildrenFiltered(":input[id$=\"__Id\"][type=hidden]")

    if len(associationIdInputs.Nodes) <= 0 {
        return nil, fmt.Errorf("cannot find option identifier element")
    }

    associationId := associationIdInputs.First().AttrOr("value", "")

    if len(associationId) <= 0 {
        return nil, fmt.Errorf("association identifier is empty")
    }

    description := associationSelection.Find("div.product-association-info div.product-association-title").Text()
    description = strings.TrimSpace(description)

    // Empty the description if it is equal to the name
    if name == description {
        description = ""
    }

    rawPrice := associationSelection.Find("div.product-association-info div.product-association-price").Text()

    var associationPrice float64

    if len(rawPrice) <= 0 {
        associationPrice = 0
    } else {
        associationPriceParts := strings.Fields(rawPrice)

        if len(associationPriceParts) != 2 {
            return nil, fmt.Errorf("invalid association price format: %s", rawPrice)
        }

        currencySymbol = associationPriceParts[0]

        optionPriceStr := strings.ReplaceAll(associationPriceParts[1], ",", ".")

        price, conversionError := strconv.ParseFloat(optionPriceStr, 64)

        if conversionError != nil {
            return nil, conversionError
        }

        associationPrice = price
    }

    association := models.ItemAssociation{
        Id:              associationId,
        GroupId:         groupId,
        Name:            name,
        Description:     description,
        IsOptional:      isOptional,
        IsAlwaysChecked: isAlwaysChecked,
        Price: models.Currency{
            Value:          associationPrice,
            CurrencySymbol: currencySymbol,
        },
    }

    return &association, nil
}
