/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/J-Siu/go-auto-docker/lib"
	"github.com/spf13/cobra"
)

// dbSearchCmd represents the dbSearch command
var dbSearchCmd = &cobra.Command{
	Use:     "search",
	Aliases: []string{"s"},
	Short:   "Search database",
	Run: func(cmd *cobra.Command, args []string) {
		lib.DbAlpine.
			Init().
			DbConnect()
		if lib.DbAlpine.Err == nil {
			for _, pkg := range args {
				lib.DbAlpine.PkgSearch(pkg)
			}
		}
	},
}

func init() {
	dbCmd.AddCommand(dbSearchCmd)
	dbSearchCmd.Flags().BoolVarP(&lib.FlagDbSearch.Exact, "exact", "e", false, "search exact word")
}
