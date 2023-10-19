package main

import (
    "encoding/json"
    "fmt"
    "os"

    "github.com/TomBarten/lunchscrape_cli/bienvenue"
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

    items := bienvenue.Scrape()

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