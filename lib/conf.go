/*
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

package lib

import (
	"github.com/J-Siu/go-helper/v2/basestruct"
	"github.com/J-Siu/go-helper/v2/errs"
	"github.com/J-Siu/go-helper/v2/ezlog"
	"github.com/J-Siu/go-helper/v2/file"
	"github.com/spf13/viper"
)

var ConfDefault = TypeConf{
	DirCache:      "~/.cache/go-auto-docker",
	DirDB:         "db",
	DirRepo:       "repo",
	FileConf:      "~/.config/go-auto-docker.json",
	FileLicense:   "LICENSE",
	FileChangeLog: "CHANGELOG.md",

	AlpineBranch: []string{"latest-stable", "edge"},

	TagReadmeLogStart: "<!--CHANGE-LOG-START-->",
	TagReadmeLogEnd:   "<!--CHANGE-LOG-END-->",
}

type TypeConf struct {
	*basestruct.Base

	DirCache      string `json:"DirCache"`    // Directory name, not full path, of cache. Default: ~/.cache/go-auto-docker
	DirDB         string `json:"DirDB"`       // Directory name, not full path, of database. Default: db
	DirRepo       string `json:"DirRepo"`     // Directory name, not full path, of repository copy. Default: repo
	FileConf      string `json:"FileConf"`    // Full path of config file. Default: ~/.config/go-auto-docker.json
	FileLicense   string `json:"FileLicense"` // Filename, not full path, of readme file. Default: LICENSE
	FileChangeLog string `json:"FileReadme"`  // Filename, not full path, of readme file. Default: README.md

	AlpineBranch []string `json:"AlpineBranch"`

	// TODO: Change following to array
	TagReadmeLogStart string `json:"ReadmeLogStart"` // Default: <!--CHANGE-LOG-START-->
	TagReadmeLogEnd   string `json:"ReadmeLogEnd"`   // Default: <!--CHANGE-LOG-END-->
}

func (t *TypeConf) New() *TypeConf {
	t.Base = new(basestruct.Base)
	t.Initialized = true
	t.MyType = "TypeConf"
	prefix := t.MyType + ".New"

	t.setDefault()
	ezlog.Debug().N(prefix).N("Default").Lm(t).Out()

	t.readFileConf()
	ezlog.Debug().N(prefix).N("Raw").Lm(t).Out()

	t.expand()
	ezlog.Debug().N(prefix).N("Expand").Lm(t).Out()

	return t
}

// viper handle file and unmarshal
func (t *TypeConf) readFileConf() *TypeConf {
	prefix := t.MyType + ".readFileConf"
	viper.SetConfigType("json")
	viper.SetConfigFile(file.TildeEnvExpand(t.FileConf))
	viper.AutomaticEnv()
	t.Err = viper.ReadInConfig()

	if t.Err == nil {
		t.Err = viper.Unmarshal(&t)
	}

	errs.Queue(prefix, t.Err)

	return t
}

// Should be called before reading config file
func (t *TypeConf) setDefault() *TypeConf {
	if t.FileConf == "" {
		t.FileConf = ConfDefault.FileConf
	}
	t.DirCache = ConfDefault.DirCache
	t.DirDB = ConfDefault.DirDB
	t.DirRepo = ConfDefault.DirRepo
	t.FileLicense = ConfDefault.FileLicense
	t.FileChangeLog = ConfDefault.FileChangeLog
	t.AlpineBranch = ConfDefault.AlpineBranch
	t.TagReadmeLogEnd = ConfDefault.TagReadmeLogEnd
	t.TagReadmeLogStart = ConfDefault.TagReadmeLogStart
	return t
}

func (t *TypeConf) expand() *TypeConf {
	t.DirCache = file.TildeEnvExpand(t.DirCache)
	t.FileConf = file.TildeEnvExpand(t.FileConf)
	return t
}
