/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/J-Siu/go-auto-docker/global"
	"github.com/J-Siu/go-helper/v2/errs"
	"github.com/spf13/cobra"
)

// dbSearchCmd represents the dbSearch command
var dbSearchCmd = &cobra.Command{
	Use:     "search",
	Aliases: []string{"s"},
	Short:   "Search database",
	Run: func(cmd *cobra.Command, args []string) {
		if global.Db.Err() == nil {
			for _, pkg := range args {
				global.Db.Search(pkg, global.FlagDbSearch.Exact)
			}
		}
		errs.Queue("", global.Db.Err())
	},
}

func init() {
	dbCmd.AddCommand(dbSearchCmd)
	dbSearchCmd.Flags().BoolVarP(&global.FlagDbSearch.Exact, "exact", "e", false, "search exact word")
}
