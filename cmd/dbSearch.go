/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

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
			var (
				strArrArr  *[]*[]string
				tab_Writer = tabwriter.NewWriter(os.Stdout, 1, 1, 1, ' ', 0)
			)
			for _, pkg := range args {
				strArrArr = global.Db.Search(pkg, global.FlagDbSearch.Exact)
				for _, strArr := range *strArrArr {
					fmt.Fprintln(tab_Writer, strings.Join(*strArr, "\t"))
				}
			}
			tab_Writer.Flush()
		}
		errs.Queue("", global.Db.Err())
	},
}

func init() {
	dbCmd.AddCommand(dbSearchCmd)
	dbSearchCmd.Flags().BoolVarP(&global.FlagDbSearch.Exact, "exact", "e", false, "search exact word")
}
