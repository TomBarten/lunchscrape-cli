package modules

import (
    "github.com/TomBarten/lunchscrape_cli/model/item"
    "github.com/TomBarten/lunchscrape_cli/modules/bienvenue"
)

type Module interface {
    Scrape() *[]item.Item
}

var ModuleMap = map[string]Module{
    "bienvenue": bienvenue.Module{},
}
