package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/TomBarten/lunchscrape_cli/bienvenue"
	"github.com/TomBarten/lunchscrape_cli/models"
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/extensions"
)

func main() {

    menuOutputFileName := "menu.json"
    menuOutputPath := "../out"

    if outputDirError := os.MkdirAll(menuOutputPath, 0755); outputDirError != nil {
        fmt.Println("Error creating directory:", outputDirError)
    }

    if outputDirError := os.Chdir(menuOutputPath); outputDirError != nil {
        fmt.Println("Error changing working directory:", outputDirError)
        return
    }

    baseDomain := "cafetariabienvenue.12waiter.eu"

    navigator := colly.NewCollector(
        colly.AllowedDomains(baseDomain),
        colly.CacheDir(fmt.Sprintf("./%s_cache", baseDomain)),
        colly.Async(true),
    )

    extensions.RandomUserAgent(navigator)

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
        bienvenue.CollectProductItem(&items, element)
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
