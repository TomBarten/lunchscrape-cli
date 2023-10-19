package collect

import (
    "fmt"
    "strconv"
    "strings"

    "github.com/TomBarten/lunchscrape_cli/models"
    "github.com/gocolly/colly"
)

func CollectItem(items *[]models.Item, element *colly.HTMLElement) {

    rawItemPrice := element.ChildText(
        "form#product-form fieldset.product-offer div.product-price-measurement div.product-price")

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

    associations := collectItemAssociations(element, currencySymbol)

    item := models.Item{
        Slug:        element.ChildAttr("input[id=\"Editor_Slug\"]", "value"),
        Name:        element.ChildText("div.product-section.product-intro h1"),
        Description: element.ChildText("div.product-section.product-intro p"),
        ImgUrl:      element.ChildAttr("div.product-image-default img", "src"),
        Price: models.Currency{
            CurrencySymbol: currencySymbol,
            Value:          itemPrice,
        },
        Options:      *options,
        Associations: *associations,
    }

    *items = append(*items, item)
}
