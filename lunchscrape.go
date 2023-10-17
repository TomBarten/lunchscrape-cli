package main

import (
    "encoding/json"
    "fmt"
    "os"
    "strconv"
    "strings"

    "github.com/PuerkitoBio/goquery"
    "github.com/gocolly/colly"
    "github.com/gocolly/colly/extensions"
)

type currency struct {
    CurrencySymbol string  `json:"currency-symbol"`
    Value          float64 `json:"value"`
}

type itemOption struct {
    Name                string   `json:"name"`
    Price               currency `json:"price"`
    GroupId             int      `json:"option-group-id"`
    IsOptional          bool     `json:"optional"`
    IsMutuallyExclusive bool     `json:"mutually-exclusive"`
}

type item struct {
    Name        string       `json:"name"`
    Price       currency     `json:"price"`
    Description string       `json:"description"`
    ImgUrl      string       `json:"img-url"`
    Options     []itemOption `json:"options"`
}

func main() {

    menuOutputFileName := "menu.json"
    menuOutputPath := "./out"

    if outputDirError := os.MkdirAll(menuOutputPath, 0755); outputDirError != nil {
        fmt.Println("Error creating directory:", outputDirError)
    }

    if outputDirError := os.Chdir(menuOutputPath); outputDirError != nil {
        fmt.Println("Error changing working directory:", outputDirError)
        return
    }

    baseDomain := "cafetariabienvenue.12waiter.eu"

    collector := colly.NewCollector(
        colly.AllowedDomains(baseDomain),
    )

    extensions.RandomUserAgent(collector)

    items := make([]item, 0, 100)

    collector.OnHTML("a.collection-item[href^='/c/']", navigate)
    collector.OnHTML("a.product-item[href^='/c/'][href*='/p/']", navigate)

    collector.OnHTML("div.product-page.product-body", func(element *colly.HTMLElement) {
        collectProductItem(&items, element)
    })

    collector.OnRequest(func(request *colly.Request) {
        fmt.Println("Visiting", request.URL)
    })

    collector.Visit(fmt.Sprintf("https://%s", baseDomain))

    jsonData, jsonError := json.Marshal(items)

    if jsonError != nil {
        fmt.Println("Error encoding JSON:", jsonError)
        return
    }

    if jsonFileError := os.WriteFile(menuOutputFileName, jsonData, 0644); jsonFileError != nil {
        fmt.Println("Error writing to file:", jsonFileError)
        return
    }
}

func navigate(element *colly.HTMLElement) {

    link := element.Attr("href")
    element.Request.Visit(link)
}

func collectItemOptions(element *colly.HTMLElement, currencySymbol string) *[]itemOption {

    optionGroups := element.DOM.Find("form#product-form div.product-section.product-section-input fieldset.product-option-group")

    options := make([]itemOption, 0, 20)

    optionGroups.Each(func(groupIteration int, optionGroupSelection *goquery.Selection) {

        isOptional := false

        optionalSpan := optionGroupSelection.Find("legend span.badge.control-optional")

        if optionalSpan != nil && len(optionalSpan.Nodes) > 0 {
            isOptional = true
        }

        handleOptionsGroup(groupIteration+1, &options, optionGroupSelection, currencySymbol, isOptional)
    })

    return &options
}

func handleOptionsGroup(groupId int, options *[]itemOption, optionGroupSelection *goquery.Selection, currencySymbol string, isOptional bool) {

    optionElements := optionGroupSelection.Find(
        "div.product-option-group-options.form-checks div.product-option.form-check")

    optionElements.Each(func(i int, optionSelection *goquery.Selection) {

        option, error := constructItemOption(
            groupId, isOptional, currencySymbol, optionSelection)

        if error != nil {
            fmt.Println(error)
            return
        }

        *options = append(*options, *option)
    })
}

func constructItemOption(
    groupId int,
    isOptional bool,
    currencySymbol string,
    optionSelection *goquery.Selection) (*itemOption, error) {

    isMutuallyExclusive := true

    checkBoxInput := optionSelection.Find("input[type=checkbox]:not([type=hidden])")

    if checkBoxInput != nil && len(checkBoxInput.Nodes) > 0 {
        isMutuallyExclusive = false
    }

    optionLabel := optionSelection.Find("label.form-check-label")

    rawOptionPrice := optionLabel.Find("span.product-option-price").Text()

    var optionPrice float64

    if len(rawOptionPrice) <= 0 {
        optionPrice = 0
    } else {
        optionPriceParts := strings.Fields(rawOptionPrice)

        if len(optionPriceParts) != 3 {
            return nil, fmt.Errorf("invalid item price format: %s", rawOptionPrice)
        }

        currencySymbol = optionPriceParts[1]

        optionPriceStr := strings.ReplaceAll(optionPriceParts[2], ",", ".")

        price, conversionError := strconv.ParseFloat(optionPriceStr, 64)

        if conversionError != nil {
            return nil, conversionError
        }

        optionPrice = price
    }

    option := itemOption{
        GroupId:             groupId,
        Name:                strings.TrimSpace(optionLabel.Contents().Not("span").Text()),
        IsOptional:          isOptional,
        IsMutuallyExclusive: isMutuallyExclusive,
        Price: currency{
            Value:          optionPrice,
            CurrencySymbol: currencySymbol,
        },
    }

    return &option, nil
}

func collectProductItem(items *[]item, element *colly.HTMLElement) {

    rawItemPrice := element.ChildText("form#product-form fieldset.product-offer div.product-price-measurement div.product-price")

    itemPriceParts := strings.Fields(rawItemPrice)

    if len(itemPriceParts) != 2 {
        fmt.Printf("invalid item price format: %s", rawItemPrice)
        return
    }

    currencySymbol := itemPriceParts[0]

    itemPriceStr := strings.ReplaceAll(itemPriceParts[1], ",", ".")

    itemPrice, conversionError := strconv.ParseFloat(itemPriceStr, 64)

    if conversionError != nil {
        fmt.Println(conversionError)
        return
    }

    options := collectItemOptions(element, currencySymbol)

    item := item{
        Name:        element.ChildText("div.product-section.product-intro h1"),
        Description: element.ChildText("div.product-section.product-intro p"),
        ImgUrl:      element.ChildAttr("div.product-image-default img", "src"),
        Price: currency{
            CurrencySymbol: currencySymbol,
            Value:          itemPrice,
        },
        Options: *options,
    }

    *items = append(*items, item)
}
