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
		global.DbAlpine.
			New(&global.Conf.DirCache, &global.Conf.DirDB, &global.Conf.AlpineBranch).
			DbConnect()
		if global.DbAlpine.Err == nil {
			for _, pkg := range args {
				global.DbAlpine.PkgSearch(pkg, global.FlagDbSearch.Exact)
			}
		}
		errs.Queue("", global.DbAlpine.Err)
	},
}

func init() {
	dbCmd.AddCommand(dbSearchCmd)
	dbSearchCmd.Flags().BoolVarP(&global.FlagDbSearch.Exact, "exact", "e", false, "search exact word")
}
