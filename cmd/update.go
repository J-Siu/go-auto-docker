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
	"github.com/J-Siu/go-auto-docker/lib"
	"github.com/J-Siu/go-helper"
	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:     "update",
	Aliases: []string{"u"},
	Short:   "Update Alpine package version",
	Run: func(cmd *cobra.Command, args []string) {
		prefix := "update"
		helper.ReportDebug(&lib.FlagUpdate, "FlagUpdate", false, false)

		var err error

		if lib.FlagUpdate.UpdateDb {
			lib.DbAlpine.
				Init().
				DbUpdate()
		} else {
			lib.DbAlpine.
				Init().
				DbConnect()
		}
		err = lib.DbAlpine.Err

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
				Init(workPath).
				CopySrcToCache()
			err = repo.Err

			// Dockerfile file
			if err == nil {
				docker.
					Init(repo.DirCache).
					Update().
					Write()
				if lib.Flag.Debug {
					docker.Dump()
				}
				err = docker.Err
			}

			// test build
			if err == nil && lib.FlagUpdate.BuildTest {
				// TODO: Docker test build
				err = docker.Err
			}

			if err == nil {
				if docker.VerNew > docker.VerCurr {

					// README.md file. Depends on docker.VerCurr. Must be done after processing docker.
					if err == nil {
						helper.Report(docker.Pkg+": "+docker.VerCurr+" -> "+docker.VerNew, prefix, false, true)
						if docker.VerNew > docker.VerCurr {
							readme.
								Init(repo.DirCache, docker.Pkg, docker.VerCurr, docker.VerNew).
								Read().
								Update().
								Write()
							if lib.Flag.Debug {
								readme.Dump()
							}
						}
						err = readme.Err
					}

					// LICENSE file
					if err == nil {
						license.
							Init(repo.DirCache).
							Read().
							Update().
							Write()
						if lib.Flag.Debug {
							license.Dump()
						}
						err = license.Err
					}

					// Repository commit and tag
					if err == nil && lib.FlagUpdate.Commit {
						repo.Commit(docker.VerNew, lib.FlagUpdate.Tag, true)
						err = repo.Err
					}

					// Repository copy back
					if err == nil && lib.FlagUpdate.Save {
						repo.CopyCacheToSrc()
					}

				} else {
					helper.Report(docker.Pkg+": "+docker.VerCurr+" -> No update", prefix, false, true)
				}
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
	updateCmd.Flags().BoolVarP(&lib.FlagUpdate.Commit, "commit", "c", false, "apply git commit. Only work with -save")
	updateCmd.Flags().BoolVarP(&lib.FlagUpdate.BuildTest, "buildTest", "b", false, "so not perform docker build")
	updateCmd.Flags().BoolVarP(&lib.FlagUpdate.Save, "save", "", false, "write back to project folder (cancel on error)")
	updateCmd.Flags().BoolVarP(&lib.FlagUpdate.Tag, "tag", "t", false, "apply git tag. (only work with --commit)")
	updateCmd.Flags().BoolVarP(&lib.FlagUpdate.UpdateDb, "updateDb", "", false, "update Alpine package database")
}
