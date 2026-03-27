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

package root

import (
	"os"

	"github.com/J-Siu/go-auto-docker/db"
	"github.com/J-Siu/go-auto-docker/global"
	"github.com/J-Siu/go-auto-docker/lib"
	"github.com/J-Siu/go-helper/v2/errs"
	"github.com/J-Siu/go-helper/v2/ezlog"
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:     "go-auto-docker",
	Version: global.Version,
	Short:   "Mass updating single package Docker project base on Alpine Linux packages.",
	Long:    `Automate update for README.md change log, apply tag according to package version. Also handle test build, git commit.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		prefix := "root"
		ezlog.SetLogLevel(ezlog.ERR)
		if global.Flag.Debug {
			ezlog.SetLogLevel(ezlog.DEBUG)
		}
		ezlog.Debug().N("Version").M(global.Version).Ln("Flag").Lm(&global.Flag).Out()
		global.Conf.New()

		global.Db = new(db.TypeDbAlpine).
			New(&global.Conf.DirCache, &global.Conf.DirDB, &global.Conf.AlpineBranch).
			Connect()
		if global.Flag.UpdateDb {
			ezlog.Log().M("db update").Out()
			global.Db.Update()
		}
		errs.Queue(prefix, global.Db.Err())
	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		if errs.NotEmpty() {
			ezlog.Err().L().M(errs.Errs()).Out()
		}
	},
}

func Execute() {
	err := RootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	RootCmd.PersistentFlags().BoolVarP(&global.Flag.Debug, "debug", "d", false, "enable debug")
	RootCmd.PersistentFlags().BoolVarP(&global.Flag.UpdateDb, "updatedb", "u", false, "update DB")
	RootCmd.PersistentFlags().BoolVarP(&global.Flag.Verbose, "verbose", "v", false, "enable verbose")
	RootCmd.PersistentFlags().StringVarP(&global.Conf.FileConf, "config", "", lib.ConfDefault.FileConf, "config file")
}
