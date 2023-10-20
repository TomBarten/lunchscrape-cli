package order

import "github.com/TomBarten/lunchscrape_cli/model"

type ItemAssociation struct {
    Id              string `json:"id"`
    Name            string `json:"name"`
    GroupId         string `json:"association-group-id"`
    IsAlwaysChecked bool   `json:"always-checked"`
}

type ItemOption struct {
    Id      string `json:"id"`
    Name    string `json:"name"`
    GroupId string `json:"option-group-id"`
}

type Item struct {
    Slug         string            `json:"slug"`
    ItemUrl      string            `json:"item-url"`
    Name         string            `json:"name"`
    Cost         model.Currency    `json:"price"`
    Options      []ItemOption      `json:"options"`
    Associations []ItemAssociation `json:"associations"`
}

type Order struct {
    Items     []Item
    TotalCost model.Currency
}
