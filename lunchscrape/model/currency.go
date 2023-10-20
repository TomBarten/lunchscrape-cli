package model

type Currency struct {
    CurrencySymbol string  `json:"currency-symbol"`
    Value          float64 `json:"value"`
}
