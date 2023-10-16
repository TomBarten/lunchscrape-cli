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

type itemOption struct {
    Name       string   `json:"name"`
    Price      currency `json:"price"`
    IsOptional bool     `json:"optional"`
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

func collectProductItemOptions(element *colly.HTMLElement) []itemOption {

    // TODO implement
    return make([]itemOption, 0)
}

func collectProductItem(items *[]item, element *colly.HTMLElement) {

    options := collectProductItemOptions(element)

    rawItemPrice := element.ChildText("form#product-form fieldset.product-offer div.product-price-measurement div.product-price")

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
        Name:        element.ChildText("div.product-section.product-intro h1"),
        Description: element.ChildText("div.product-section.product-intro p"),
        ImgUrl:      element.ChildAttr("div.product-image-default img", "src"),
        Price: currency{
            CurrencySymbol: currencySymbol,
            Value:          itemPrice,
        },
        Options: options,
    }

    *items = append(*items, item)
}
