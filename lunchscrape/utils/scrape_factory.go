package utils

import (
    "fmt"

    "github.com/gocolly/colly"
    "github.com/gocolly/colly/extensions"
)

func CreateScraper(baseDomain string, withCache bool) *colly.Collector {

    collector := colly.NewCollector(
        colly.AllowedDomains(baseDomain),
        colly.Async(true),
    )

    extensions.RandomUserAgent(collector)

    if withCache {
        collector.CacheDir = fmt.Sprintf("./%s_cache", baseDomain)
    }

    return collector
}
