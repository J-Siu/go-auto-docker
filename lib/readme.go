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

package lib

import (
	"errors"
	"os"
	"path"
	"regexp"
	"strings"
	"time"

	"github.com/J-Siu/go-helper"
)

type TypeReadme struct {
	Err    error
	init   bool
	myType string

	Content  []string
	Dir      string
	FilePath string

	TagReadmeLogStart string
	TagReadmeLogEnd   string

	Pkg     string
	VerCurr string
	VerNew  string

	changeLogStart bool
	changeLogEnd   bool
	newLogHeader   string
	newLogContent  string
}

func (r *TypeReadme) Init(dir string, pkg string, verCurr string, verNew string) *TypeReadme {
	r.init = true
	r.myType = "TypeReadme"
	r.Dir = dir
	r.FilePath = path.Join(r.Dir, Conf.FileReadme)
	r.TagReadmeLogStart = Conf.TagReadmeLogStart
	r.TagReadmeLogEnd = Conf.TagReadmeLogEnd

	r.Pkg = pkg
	r.VerCurr = verCurr
	r.VerNew = verNew

	prefix := r.myType + ".Init"
	helper.ReportDebug(r, prefix, false, true)
	if !helper.IsRegularFile(r.FilePath) {
		r.Err = errors.New("TypeReadme.Init: " + r.FilePath + " not found")
		helper.ErrsQueue(r.Err, prefix)
	}
	return r
}

func (r *TypeReadme) Dump() *TypeReadme {
	prefix := r.myType + ".Dump"
	if r.Err != nil {
		return r
	}
	if !r.init {
		r.Err = errors.New("not initialized")
		helper.ErrsQueue(r.Err, prefix)
	}
	if r.Err == nil {
		helper.Report(r, prefix, false, false)
	}
	return r
}

func (r *TypeReadme) Update() *TypeReadme {
	prefix := r.myType + ".Update"
	if r.Err != nil {
		return r
	}
	if !r.init {
		r.Err = errors.New("not initialized")
		helper.ErrsQueue(r.Err, prefix)
	}
	if r.Err == nil {
		if r.VerNew > r.VerCurr {
			helper.ReportDebug(r.Pkg+": "+r.VerCurr+" -> "+r.VerNew, prefix, false, true)
			r.changeLogStart = false
			r.changeLogEnd = false
			r.newLogHeader = "- " + r.VerNew
			r.newLogContent = "  - Auto update to " + r.VerNew
			// for lineNum := range r.Content {
			for lineNum := range r.Content {
				if r.Err == nil {
					helper.ReportDebug(&r.Content[lineNum], prefix, false, true)
					r.
						updateLog(&r.Content[lineNum]).
						updateLicenseYear(&r.Content[lineNum])
					helper.ReportDebug(&r.Content[lineNum], prefix, false, true)
				}
			}
		}
	}
	return r
}

// Read README.md into `Content`
func (r *TypeReadme) Read() *TypeReadme {
	prefix := r.myType + ".Read"
	if r.Err != nil {
		return r
	}
	if !r.init {
		r.Err = errors.New("not initialized")
	}
	if r.Err == nil {
		r.Content, r.Err = helper.FileStrArrRead(r.FilePath)
		if r.Err != nil {
			r.Err = errors.New(r.FilePath + " not found")
		} else {
		}
	}
	helper.ErrsQueue(r.Err, prefix)
	return r
}

// Write `Content` into README.md
func (r *TypeReadme) Write() *TypeReadme {
	prefix := r.myType + ".Write"
	if r.Err != nil {
		return r
	}
	if !r.init {
		r.Err = errors.New("not initialized")
	}
	if r.Err == nil {
		fileStats, err := os.Stat(r.FilePath)
		if err != nil {
			// Should never happen at this stage, but ...
			r.Err = err
		} else {
			helper.FileStrArrWrite(r.FilePath, r.Content, fileStats.Mode())
		}
	}
	helper.ErrsQueue(r.Err, prefix)
	return r
}

func (r *TypeReadme) updateLog(line *string) *TypeReadme {
	prefix := r.myType + ".updateLog"
	if strings.EqualFold(*line, r.TagReadmeLogStart) {
		r.changeLogStart = true
		helper.ReportDebug("changeStart", prefix, false, true)
	}
	if !r.changeLogEnd && r.changeLogStart {
		helper.ReportDebug("changeStarted "+*line, prefix, false, true)
		if strings.Contains(*line, r.newLogHeader) {
			// Something wrong if new log header already exist
			r.Err = errors.New("change log already contains " + r.VerNew)
			helper.ErrsQueue(r.Err, prefix)
		}
		if r.Err == nil && strings.Contains(*line, r.TagReadmeLogEnd) {
			*line = r.newLogHeader + "\n" + r.newLogContent + "\n" + *line
			helper.ReportDebug(line, prefix, false, false)
		}
	}
	if strings.Contains(*line, r.TagReadmeLogEnd) {
		r.changeLogEnd = true
		helper.ReportDebug("changeEnd", prefix, false, true)
	}
	helper.ReportDebug(line, prefix, false, true)
	return r
}

// This work for:
//   - Copyright (c) xxxx
//   - Copyright © xxxx
func (r *TypeReadme) updateLicenseYear(line *string) *TypeReadme {
	prefix := r.myType + ".updateLicenseYear"
	c := []string{"(c)", "©"}
	cRaw := []string{`\(c\)`, `©`}
	copyright := "Copyright"
	year, _ := helper.NumToStr(time.Now().Year())
	for i := range c {
		// "(?i)" <- means case insensitive
		re := regexp.MustCompile("(?i)" + copyright + " " + cRaw[i] + ` \d\d\d\d`)
		*line = re.ReplaceAllString(*line, copyright+" "+c[i]+" "+year)
	}
	helper.ReportDebug(line, prefix, false, true)
	return r
}
