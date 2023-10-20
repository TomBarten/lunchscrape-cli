package item

import "github.com/TomBarten/lunchscrape_cli/model"

type ItemOption struct {
    Id                  string         `json:"id"`
    Name                string         `json:"name"`
    Price               model.Currency `json:"price"`
    GroupId             string         `json:"option-group-id"`
    IsOptional          bool           `json:"optional"`
    IsMutuallyExclusive bool           `json:"mutually-exclusive"`
}

type ItemAssociation struct {
    Id              string         `json:"id"`
    Name            string         `json:"name"`
    Description     string         `json:"description"`
    Price           model.Currency `json:"price"`
    GroupId         string         `json:"association-group-id"`
    IsAlwaysChecked bool           `json:"always-checked"`
    IsOptional      bool           `json:"optional"`
}

type Item struct {
    Slug         string            `json:"slug"`
    Name         string            `json:"name"`
    Price        model.Currency    `json:"price"`
    Description  string            `json:"description"`
    ItemUrl      string            `json:"item-url"`
    ImgUrl       string            `json:"img-url"`
    Options      []ItemOption      `json:"options"`
    Associations []ItemAssociation `json:"associations"`
}
