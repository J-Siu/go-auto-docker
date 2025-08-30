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
	"github.com/J-Siu/go-helper"
	"github.com/spf13/viper"
)

var Default = TypeConf{
	DirCache:    "~/.cache/go-auto-docker",
	DirDB:       "db",
	DirProj:     "proj",
	FileConf:    "~/.config/go-auto-docker.json",
	FileLicense: "LICENSE",
	FileReadme:  "README.md",

	TagReadmeLogStart: "<!--CHANGE-LOG-START-->",
	TagReadmeLogEnd:   "<!--CHANGE-LOG-END-->",
}

type TypeConf struct {
	Err    error
	myType string
	init   bool

	DirCache    string `json:"DirCache"`    // Directory name, not full path, of cache. Default: ~/.cache/go-auto-docker
	DirDB       string `json:"DirDB"`       // Directory name, not full path, of database. Default: db
	DirProj     string `json:"DirProj"`     // Directory name, not full path, of project(cloning). Default: proj
	FileConf    string `json:"FileConf"`    // Full path of config file. Default: ~/.config/go-auto-docker.json
	FileLicense string `json:"FileLicense"` // Filename, not full path, of readme file. Default: LICENSE
	FileReadme  string `json:"FileReadme"`  // Filename, not full path, of readme file. Default: README.md

	TagReadmeLogStart string `json:"ReadmeLogStart"` // Default: <!--CHANGE-LOG-START-->
	TagReadmeLogEnd   string `json:"ReadmeLogEnd"`   // Default: <!--CHANGE-LOG-END-->
}

func (c *TypeConf) Init() *TypeConf {
	c.init = true
	c.myType = "TypeConf"
	prefix := c.myType + ".Init"

	c.readFileConf()
	helper.ReportDebug(c, prefix+": Raw", false, true)

	c.setDefault()
	helper.ReportDebug(c, prefix+": Default + Flag", false, true)

	c.expand()
	helper.ReportDebug(c, prefix+": Expand", false, true)

	return c
}

// viper handle file and unmarshal
func (c *TypeConf) readFileConf() *TypeConf {
	prefix := c.myType + ".readFileConf"
	viper.SetConfigType("json")
	viper.SetConfigFile(helper.TildeEnvExpand(Conf.FileConf))
	viper.AutomaticEnv()
	c.Err = viper.ReadInConfig()

	if c.Err == nil {
		c.Err = viper.Unmarshal(&c)
	}

	helper.ErrsQueue(c.Err, prefix)

	return c
}

// Set default value if a field is empty
func (c *TypeConf) setDefault() *TypeConf {
	if c.DirCache == "" {
		c.DirCache = Default.DirCache
	}
	if c.DirDB == "" {
		c.DirDB = Default.DirDB
	}
	if c.DirProj == "" {
		c.DirProj = Default.DirProj
	}
	if c.FileConf == "" {
		c.FileConf = Default.FileConf
	}
	if c.FileLicense == "" {
		c.FileLicense = Default.FileLicense
	}
	if c.FileReadme == "" {
		c.FileReadme = Default.FileReadme
	}
	if c.TagReadmeLogEnd == "" {
		c.TagReadmeLogEnd = Default.TagReadmeLogEnd
	}
	if c.TagReadmeLogStart == "" {
		c.TagReadmeLogStart = Default.TagReadmeLogStart
	}
	return c
}

func (c *TypeConf) expand() *TypeConf {
	c.DirCache = helper.TildeEnvExpand(c.DirCache)
	c.FileConf = helper.TildeEnvExpand(c.FileConf)
	return c
}
