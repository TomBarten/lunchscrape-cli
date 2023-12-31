package collect

import (
    "fmt"
    "strconv"
    "strings"

    "github.com/PuerkitoBio/goquery"
    "github.com/TomBarten/lunchscrape_cli/model"
    "github.com/TomBarten/lunchscrape_cli/model/item"
    "github.com/gocolly/colly"
)

func collectItemOptions(element *colly.HTMLElement, currencySymbol string) *[]item.ItemOption {

    optionGroups := element.DOM.Find(
        "form#product-form div.product-section.product-section-input fieldset.product-option-group")

    if len(optionGroups.Nodes) <= 0 {
        empty := make([]item.ItemOption, 0)
        return &empty
    }

    options := make([]item.ItemOption, 0, 20)

    optionGroups.EachWithBreak(func(groupIteration int, groupSelection *goquery.Selection) bool {

        isOptional := false

        optionalSpan := groupSelection.Find("legend span.badge.control-optional")

        if optionalSpan != nil && len(optionalSpan.Nodes) > 0 {
            isOptional = true
        }

        error := handleOptionsGroup(&options, groupSelection, currencySymbol, isOptional)

        if error == nil {
            return true
        }

        fmt.Println(error)
        return false
    })

    return &options
}

func handleOptionsGroup(
    options *[]item.ItemOption,
    optionGroupSelection *goquery.Selection,
    currencySymbol string,
    isOptional bool) error {

    optionElements := optionGroupSelection.Find(
        "div.product-option-group-options.form-checks div.product-option.form-check")

    optionGroupIdInputs := optionGroupSelection.ChildrenFiltered(":input[id$=\"__Id\"][type=hidden]")

    if len(optionGroupIdInputs.Nodes) <= 0 {
        return fmt.Errorf("cannot find option group identifier element")
    }

    groupId := optionGroupIdInputs.First().AttrOr("value", "")

    if len(groupId) <= 0 {
        return fmt.Errorf("option group identifier is empty")
    }

    optionElements.EachWithBreak(func(i int, optionSelection *goquery.Selection) bool {

        option, error := constructItemOption(
            groupId, isOptional, currencySymbol, optionSelection)

        if error != nil {
            fmt.Println(error)
            return false
        }

        *options = append(*options, *option)
        return true
    })

    return nil
}

func constructItemOption(
    groupId string,
    isOptional bool,
    currencySymbol string,
    optionSelection *goquery.Selection) (*item.ItemOption, error) {

    isMutuallyExclusive := true

    checkBoxInput := optionSelection.Find("input[type=checkbox]:not([type=hidden])")

    if checkBoxInput != nil && len(checkBoxInput.Nodes) > 0 {
        isMutuallyExclusive = false
    }

    optionIdInputs := optionSelection.ChildrenFiltered(":input[id$=\"__Id\"][type=hidden]")

    if len(optionIdInputs.Nodes) <= 0 {
        return nil, fmt.Errorf("cannot find option identifier element")
    }

    optionId := optionIdInputs.First().AttrOr("value", "")

    if len(optionId) <= 0 {
        return nil, fmt.Errorf("option identifier is empty")
    }

    optionLabel := optionSelection.Find("label.form-check-label")

    rawOptionPrice := optionLabel.Find("span.product-option-price").Text()

    var optionPrice float64

    if len(rawOptionPrice) <= 0 {
        optionPrice = 0
    } else {
        optionPriceParts := strings.Fields(rawOptionPrice)

        if len(optionPriceParts) != 3 {
            return nil, fmt.Errorf("invalid option price format: %s", rawOptionPrice)
        }

        currencySymbol = optionPriceParts[1]

        optionPriceStr := strings.ReplaceAll(optionPriceParts[2], ",", ".")

        price, conversionError := strconv.ParseFloat(optionPriceStr, 64)

        if conversionError != nil {
            return nil, conversionError
        }

        optionPrice = price
    }

    option := item.ItemOption{
        Id:                  optionId,
        GroupId:             groupId,
        Name:                strings.TrimSpace(optionLabel.Contents().Not("span").Text()),
        IsOptional:          isOptional,
        IsMutuallyExclusive: isMutuallyExclusive,
        Price: model.Currency{
            Value:          optionPrice,
            CurrencySymbol: currencySymbol,
        },
    }

    return &option, nil
}
