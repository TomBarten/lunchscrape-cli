package bienvenue

import (
    "fmt"

    "github.com/TomBarten/lunchscrape_cli/model/item"
    "github.com/TomBarten/lunchscrape_cli/modules/bienvenue/internal/collect"
    "github.com/TomBarten/lunchscrape_cli/utils"
    "github.com/gocolly/colly"
    "github.com/gocolly/colly/extensions"
)

type Module struct{}

func (m Module) Scrape() *[]item.Item {

    baseDomain := "cafetariabienvenue.12waiter.eu"

    navigator := utils.CreateScraper(baseDomain, true)

    dataCollector := navigator.Clone()

    extensions.RandomUserAgent(dataCollector)

    items := make([]item.Item, 0, 250)

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
        collect.CollectItem(&items, element)
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
