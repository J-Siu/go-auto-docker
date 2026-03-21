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
	"github.com/J-Siu/go-auto-docker/global"
	"github.com/J-Siu/go-auto-docker/lib"
	"github.com/J-Siu/go-helper/v2/errs"
	"github.com/J-Siu/go-helper/v2/ezlog"
	"github.com/spf13/cobra"
)

// checkCmd represents the check command
var checkCmd = &cobra.Command{
	Use:     "check",
	Aliases: []string{"c"},
	Short:   "Check Alpine package update",
	Run: func(cmd *cobra.Command, args []string) {
		prefix := "Update"

		var (
			err error
		)

		if err != nil {
			return
		}

		if len(args) == 0 {
			args = []string{"."}
		}

		for _, workPath := range args {
			docker := lib.TypeDocker{}
			repo := lib.TypeRepository{}

			// Repository copy to cache(tmp)
			repo.
				New(&workPath, &global.Conf.DirCache, &global.Conf.DirRepo, global.Flag.Verbose).
				CopySrcToCache()
			err = repo.Err

			// Dockerfile file
			if err == nil {
				docker.
					New(&repo.DirSrc, global.Db, global.Flag.Debug, global.Flag.Verbose).
					Check()
				err = docker.Err
			}

			if err == nil {
				ezlog.Log().N(prefix).N("YES").N(docker.Pkg).M(docker.VerCurr).M("->").M(docker.VerNew).Out()
			}

			errs.Queue("", err)
		}
	},
}

func init() {
	cmd := checkCmd
	rootCmd.AddCommand(cmd)
}
