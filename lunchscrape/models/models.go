package models

type Currency struct {
    CurrencySymbol string  `json:"currency-symbol"`
    Value          float64 `json:"value"`
}

type ItemOption struct {
    Id                  string   `json:"id"`
    Name                string   `json:"name"`
    Price               Currency `json:"price"`
    GroupId             string   `json:"option-group-id"`
    IsOptional          bool     `json:"optional"`
    IsMutuallyExclusive bool     `json:"mutually-exclusive"`
}

type ItemAssociation struct {
    Id         string   `json:"id"`
    Name       string   `json:"name"`
    Price      Currency `json:"price"`
    GroupId    string   `json:"option-group-id"`
    IsOptional bool     `json:"optional"`
}

type Item struct {
    Slug         string            `json:"slug"`
    Name         string            `json:"name"`
    Price        Currency          `json:"price"`
    Description  string            `json:"description"`
    ImgUrl       string            `json:"img-url"`
    Options      []ItemOption      `json:"options"`
    Associations []ItemAssociation `json:"associations"`
}
