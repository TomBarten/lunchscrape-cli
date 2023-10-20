package order

import (
    "encoding/json"

    "github.com/TomBarten/lunchscrape_cli/model"
    "github.com/TomBarten/lunchscrape_cli/model/item"
)

func DeserializeJsonOrder(rawItems []byte) (*Order, error) {

    var items []item.Item

    jsonDeserializeError := json.Unmarshal(rawItems, &items)

    if jsonDeserializeError != nil {
        return nil, jsonDeserializeError
    }

    order, err := transformIntoOrderItems(&items)

    if err != nil {
        return nil, err
    }

    return order, nil
}

func transformIntoOrderItems(items *[]item.Item) (*Order, error) {

    var totalCostValue float64

    orderItems := make([]Item, 0, len(*items))

    var currencySymbol string

    for _, item := range *items {

        totalCostValue += item.Price.Value

        // Only have to set once
        if len(currencySymbol) <= 0 {
            currencySymbol = item.Price.CurrencySymbol
        }

        itemJson, jsonDeserializeError := json.Marshal(item)

        if jsonDeserializeError != nil {
            return nil, jsonDeserializeError
        }

        var orderItem Item

        jsonDeserializeError = json.Unmarshal(itemJson, &orderItem)

        if jsonDeserializeError != nil {
            return nil, jsonDeserializeError
        }

        orderItems = append(orderItems, orderItem)

        for _, option := range item.Options {

            totalCostValue += option.Price.Value
        }

        for _, association := range item.Associations {

            totalCostValue += association.Price.Value
        }
    }

    order := Order{
        Items: orderItems,
        TotalCost: model.Currency{
            Value:          totalCostValue,
            CurrencySymbol: currencySymbol,
        },
    }

    return &order, nil
}
