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
	"github.com/J-Siu/go-helper/v2/ezlog"
	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:     "update",
	Aliases: []string{"u"},
	Short:   "Update Alpine package version",
	Run: func(cmd *cobra.Command, args []string) {
		prefix := "update"
		ezlog.Debug().Nn("FlagUpdate").M(&global.FlagUpdate).Out()

		var err error

		if global.FlagUpdate.UpdateDb {
			global.DbAlpine.
				New(&global.Conf.DirCache, &global.Conf.DirDB, &global.Conf.AlpineBranch).
				DbUpdate()
		} else {
			global.DbAlpine.
				New(&global.Conf.DirCache, &global.Conf.DirDB, &global.Conf.AlpineBranch).
				DbConnect()
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
			license := lib.TypeLicense{}
			repo := lib.TypeRepository{}
			readme := lib.TypeReadme{}

			// Repository copy to cache(tmp)
			repo.
				New(&workPath, &global.Conf.DirCache, &global.Conf.DirRepo, global.Flag.Verbose).
				CopySrcToCache()
			err = repo.Err

			// Dockerfile file
			if err == nil {
				docker.
					New(&repo.DirCache, global.Flag.Debug, global.Flag.Verbose).
					Update(&global.DbAlpine)
				if global.Flag.Debug {
					docker.Dump()
				}
				err = docker.Err
			}

			if err == nil {
				if docker.VerNew > docker.VerCurr {

					if err == nil {
						docker.Write()
						err = docker.Err
					}

					// test build
					if err == nil && global.FlagUpdate.BuildTest {
						docker.BuildTest()
						err = docker.Err
					}

					// README.md file. Depends on docker.VerCurr. Must be done after processing docker.
					if err == nil {
						readme.
							New(&repo.DirCache, &docker.Pkg, &docker.VerCurr, &docker.VerNew, &global.Conf.FileReadme, &global.Conf.TagReadmeLogStart, &global.Conf.TagReadmeLogEnd).
							Read().
							Update().
							Write()
						if global.Flag.Debug {
							readme.Dump()
						}
						err = readme.Err
					}

					// LICENSE file
					if err == nil {
						license.
							New(&repo.DirCache, &global.Conf.FileLicense).
							Read().
							Update().
							Write()
						if global.Flag.Debug {
							license.Dump()
						}
						err = license.Err
					}

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
					verNew := docker.VerNew
					if verNew == "" {
						verNew = "N/A"
					}
					ezlog.Log().N(prefix).N("NO").N(docker.Pkg).M(docker.VerCurr).M("->").M(docker.VerNew).Out()
				}
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
	updateCmd.Flags().BoolVarP(&global.FlagUpdate.Commit, "commit", "c", false, "apply git commit. Only work with -save")
	updateCmd.Flags().BoolVarP(&global.FlagUpdate.BuildTest, "buildTest", "b", false, "so not perform docker build")
	updateCmd.Flags().BoolVarP(&global.FlagUpdate.Save, "save", "s", false, "write back to project folder (cancel on error)")
	updateCmd.Flags().BoolVarP(&global.FlagUpdate.Tag, "tag", "t", false, "apply git tag. (only work with --commit)")
	updateCmd.Flags().BoolVarP(&global.FlagUpdate.UpdateDb, "updateDb", "u", false, "update Alpine package database")
}
