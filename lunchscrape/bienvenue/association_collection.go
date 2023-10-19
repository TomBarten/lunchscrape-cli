package bienvenue

import (
    "github.com/TomBarten/lunchscrape_cli/models"
    "github.com/gocolly/colly"
)

func collectItemAssociations(element *colly.HTMLElement, currencySymbol string) *[]models.ItemAssociation {
    associations := make([]models.ItemAssociation, 0, 10)

    return &associations
}
