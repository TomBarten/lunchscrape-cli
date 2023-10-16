package main

import (
    "encoding/json"
    "fmt"
    "os"
    "strconv"
    "strings"

    "github.com/gocolly/colly"
    "github.com/gocolly/colly/extensions"
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

    dirPath := "./out"

    if outputDirError := os.MkdirAll(dirPath, 0755); outputDirError != nil {
        fmt.Println("Error creating directory:", outputDirError)
    }

    if outputDirError := os.Chdir(dirPath); outputDirError != nil {
        fmt.Println("Error changing working directory:", outputDirError)
        return
    }

    baseDomain := "cafetariabienvenue.12waiter.eu"

    collector := colly.NewCollector(
        colly.AllowedDomains(baseDomain),
    )

    extensions.RandomUserAgent(collector)

    items := make([]item, 0, 500)

    collector.OnHTML("a.collection-item[href^='/c/']", navigate)

    collector.OnHTML("a[class=product-item]", func(element *colly.HTMLElement) {
        collectProductItems(&items, element)
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

    jsonFileError := os.WriteFile("menu.json", jsonData, 0644)

    if jsonFileError != nil {
        fmt.Println("Error writing to file:", jsonFileError)
        return
    }
}

func navigate(element *colly.HTMLElement) {

    link := element.Attr("href")
    element.Request.Visit(link)
}

func collectProductItemOptions(element *colly.HTMLElement) {
    navigate(element)
}

func collectProductItems(items *[]item, element *colly.HTMLElement) {

    collectProductItemOptions(element)

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

    *items = append(*items, item)
}
