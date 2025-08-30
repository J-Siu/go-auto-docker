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
	"fmt"

	"github.com/J-Siu/go-helper"
	"github.com/go-git/go-git/v6"
	"github.com/spf13/cobra"
)

// patchCmd represents the patch command
var patchCmd = &cobra.Command{
	Use:     "patch",
	Aliases: []string{"p"},
	Short:   "Add a patch level",
	Long:    `Add a patch level after alpine version.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("patch called")
		if len(args) == 0 {
			args = []string{"."}
		}
		for _, workPath := range args {
			helper.Report(&workPath, "", false, true)
			if _, err := git.PlainOpen(workPath); err != nil {
				helper.Report("is not a git repository.", workPath, true, true)
				continue
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(patchCmd)
}
