package main

import (
    "encoding/json"
    "fmt"
    "os"

    "github.com/TomBarten/lunchscrape_cli/cmd"
    "github.com/TomBarten/lunchscrape_cli/modules/bienvenue"
)

func main() {

    cmd.Execute()

    menuOutputFileName := "menu.json"
    outputPath := "../output"

    if outputDirError := os.MkdirAll(outputPath, 0755); outputDirError != nil {
        fmt.Println("Error creating directory:", outputDirError)
    }

    if outputDirError := os.Chdir(outputPath); outputDirError != nil {
        fmt.Println("Error changing working directory:", outputDirError)
        return
    }

    return

    items := bienvenue.Module{}.Scrape()

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
