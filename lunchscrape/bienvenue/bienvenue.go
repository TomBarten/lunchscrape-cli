package bienvenue

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/TomBarten/lunchscrape_cli/models"
	"github.com/TomBarten/lunchscrape_cli/utils"
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/extensions"
)

func Scrape() *[]models.Item {

    baseDomain := "cafetariabienvenue.12waiter.eu"

    navigator := utils.CreateScraper(baseDomain, true)

    dataCollector := navigator.Clone()

    extensions.RandomUserAgent(dataCollector)

    items := make([]models.Item, 0, 250)

    navigator.OnHTML("a.collection-item[href^='/c/']", func(element *colly.HTMLElement) {
        link := element.Attr("href")
        element.Request.Visit(link)
    })

    navigator.OnHTML("a.product-item[href^='/c/'][href*='/p/']", func(element *colly.HTMLElement) {
        productUrl := element.Request.AbsoluteURL(element.Attr("href"))
        if len(productUrl) > 0 {
            dataCollector.Visit(productUrl)
        }
    })

    dataCollector.OnHTML("div.product-page.product-body", func(element *colly.HTMLElement) {
        collectProductItem(&items, element)
    })

    navigator.OnRequest(func(request *colly.Request) {
        fmt.Println("Navigator, visiting:", request.URL)
    })

    dataCollector.OnRequest(func(request *colly.Request) {
        fmt.Println("Data collector, visiting:", request.URL)
    })

    navigator.Visit(fmt.Sprintf("https://%s", baseDomain))

    navigator.Wait()
    dataCollector.Wait()

    return &items
}

func collectProductItem(items *[]models.Item, element *colly.HTMLElement) {

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
