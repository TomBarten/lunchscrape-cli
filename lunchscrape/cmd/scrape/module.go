/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package scrape

import (
	"fmt"

	"github.com/TomBarten/lunchscrape_cli/modules"
	"github.com/spf13/cobra"
)

// moduleCmd represents the module command
var moduleCmd = &cobra.Command{
	Use:   "module",
	Short: "Module to use to scrape",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		if listModules {
			moduleNamesList()
		}
	},
}

var (
	moduleName  string
	listModules bool
)

func moduleNamesList() {

	fmt.Println("Available modules:")
	for moduleName := range modules.ModuleMap {
		fmt.Println("-", moduleName)
	}
}

func init() {

	moduleCmd.Flags().StringVarP(&moduleName, "name", "n", "", "The module name")
	moduleCmd.Flags().BoolVarP(&listModules, "list", "l", false, "Lists the available modules")

	moduleCmd.MarkFlagsMutuallyExclusive("name", "list")

	ScrapeCmd.AddCommand(moduleCmd)

	moduleCmd.PreRun = func(cmd *cobra.Command, args []string) {

	}

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// moduleCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// moduleCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
