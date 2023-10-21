# Scraping lunch/food related websites
Hobby project which can scrape 'shop' data from food-related sites.

## Development
#### [Cobra CLI](https://github.com/spf13/cobra)
```bash
go install github.com/spf13/cobra-cli@latest
```

## Objectives
- [x] Scrape single statically configured page
- [x] Scrape multiple statically configured pages
- [ ] Make scraping configurable through variables
- [ ] Make scraping configurable through command line arguments
- [x] Make scraping data outputted as JSON
    - Outputted in the `out` folder as `menu.json` 