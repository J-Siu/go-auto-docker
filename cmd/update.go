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

var updateCmd = &cobra.Command{
	Use:     "update",
	Aliases: []string{"u"},
	Short:   "Update Alpine package version",
	Run: func(cmd *cobra.Command, args []string) {
		prefix := "update"
		ezlog.Debug().N("FlagUpdate").Lm(&global.FlagUpdate).Out()

		var (
			err error
		)

		global.DbAlpine.New(&global.Conf.DirCache, &global.Conf.DirDB, &global.Conf.AlpineBranch)
		if global.FlagUpdate.UpdateDb {
			global.DbAlpine.DbUpdate()
		} else {
			global.DbAlpine.DbConnect()
		}
		err = global.DbAlpine.Err

		if err != nil {
			return
		}

		if len(args) == 0 {
			args = []string{"."}
		}

		for _, workPath := range args {
			docker := lib.TypeDocker{}
			repo := lib.TypeRepository{}
			changelog := lib.TypeChangeLog{}

			// Repository copy to cache(tmp)
			repo.
				New(&workPath, &global.Conf.DirCache, &global.Conf.DirRepo, global.Flag.Verbose).
				CopySrcToCache()
			err = repo.Err

			// Dockerfile file
			if err == nil {
				docker.
					New(&repo.DirCache, global.Flag.Debug, global.Flag.Verbose).
					Update(&global.DbAlpine).
					Write().
					Dump(global.Flag.Debug).
					BuildTest(global.FlagUpdate.BuildTest)
				err = docker.Err
			}

			if err == nil {
				if docker.Updated() {
					// README.md file. Depends on docker.VerCurr. Must be done after processing docker.
					property := lib.TypeChangeLogProperty{
						Dir:           &repo.DirCache,
						FileChangeLog: &global.Conf.FileChangeLog,
						Pkg:           &docker.Pkg,
						VerCurr:       &docker.VerCurr,
						VerNew:        &docker.VerNew,
					}
					changelog.
						New(&property).
						Update().
						Write().
						Dump(global.Flag.Debug)
					err = changelog.Err

					// Repository commit and tag
					if err == nil && global.FlagUpdate.Commit {
						repo.Commit(docker.VerNew, global.FlagUpdate.Tag, true)
						err = repo.Err
					}

					// Repository copy back
					if err == nil && global.FlagUpdate.Save {
						repo.CopyCacheToSrc()
					}

					if err == nil {
						ezlog.Log().N(prefix).N("YES").N(docker.Pkg).M(docker.VerCurr).M("->").M(docker.VerNew).Out()
					}
				} else {
					ezlog.Log().N(prefix).N("NO").N(docker.Pkg).M(docker.VerCurr).M("->")
					if docker.VerNew == "" {
						ezlog.M("<package not found>")
					} else {
						ezlog.M(docker.VerNew)
					}
					ezlog.Out()
				}
			}
			errs.Queue("", err)
		}
	},
}

func init() {
	cmd := updateCmd
	rootCmd.AddCommand(cmd)
	cmd.Flags().BoolVarP(&global.FlagUpdate.Commit, "commit", "c", false, "apply git commit. Only work with -save")
	cmd.Flags().BoolVarP(&global.FlagUpdate.BuildTest, "buildTest", "b", false, "so not perform docker build")
	cmd.Flags().BoolVarP(&global.FlagUpdate.Save, "save", "s", false, "write back to project folder (cancel on error)")
	cmd.Flags().BoolVarP(&global.FlagUpdate.Tag, "tag", "t", false, "apply git tag. (only work with --commit)")
	cmd.Flags().BoolVarP(&global.FlagUpdate.UpdateDb, "updateDb", "u", false, "update Alpine package database")
}
