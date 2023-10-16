package main

import (
    "fmt"
    "strconv"
    "strings"

    "github.com/gocolly/colly"
)

type currency struct {
    CurrencySymbol string  `json:"currency-symbol"`
    Value          float64 `json:"value"`
}

type item struct {
    Name        string `json:"name"`
    Price       currency
    Description string `json:"description"`
    ImgUrl      string `json:"imgurl"`
}

func main() {

    baseDomain := "cafetariabienvenue.12waiter.eu"

    collector := colly.NewCollector(
        colly.AllowedDomains(baseDomain),
    )

    collector.OnHTML("a.collection-item[href^='/c/']", func(h *colly.HTMLElement) {

    })

    collector.OnHTML("a[class=product-item]", collectProductItems)

    collector.OnRequest(func(request *colly.Request) {
        fmt.Println("Visiting", request.URL)
    })

    collector.Visit(fmt.Sprintf("https://%s/c/topped-friet", baseDomain))
}

func collectProductItems(element *colly.HTMLElement) {
    rawItemPrice := element.ChildText("div.product-item-body div.product-item-offer span.product-item-price")

    itemPriceParts := strings.Fields(rawItemPrice)

    if len(itemPriceParts) != 2 {
        fmt.Printf("Invalid item price format: %s", rawItemPrice)
        return
    }

    currencySymbol := itemPriceParts[0]

    itemPriceStr := strings.ReplaceAll(itemPriceParts[1], ",", ".")

    itemPrice, conversionError := strconv.ParseFloat(itemPriceStr, 64)

    if conversionError != nil {
        fmt.Println(conversionError)
        return
    }

    item := item{
        Name:        element.ChildText("div.product-item-body div[class=product-item-title]"),
        Description: element.ChildText("div.product-item-body div[class=product-item-description]"),
        Price: currency{
            CurrencySymbol: currencySymbol,
            Value:          itemPrice,
        },
    }
    fmt.Println(item)
}
