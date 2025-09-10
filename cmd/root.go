/*
The MIT License (MIT)

Copyright © 2025 John, Sing Dao, Siu <john.sd.siu@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/

package cmd

import (
	"os"

	"github.com/J-Siu/go-auto-docker/lib"

	"github.com/J-Siu/go-helper"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     "go-auto-docker",
	Version: lib.Version,
	Short:   "Mass updating single package Docker project base on Alpine Linux packages.",
	Long:    `Automate update for README.md change log, apply tag according to package version. Also handle test build, git commit.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		helper.Debug = lib.Flag.Debug
		helper.ReportDebug(lib.Version, "Version", false, true)
		lib.Conf.Init()
		helper.ReportDebug(&lib.Flag, "Flag", false, false)
	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		helper.Report(helper.Errs, "Errors", true, false)
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&lib.Flag.Debug, "debug", "d", false, "enable debug")
	rootCmd.PersistentFlags().BoolVarP(&lib.Flag.Verbose, "verbose", "v", false, "enable debug")
	rootCmd.PersistentFlags().StringVarP(&lib.Conf.FileConf, "config", "", lib.ConfDefault.FileConf, "config file")
}
