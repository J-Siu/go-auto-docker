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
	"time"

	"github.com/J-Siu/go-helper"
)

type TypeLicense struct {
	Err    error
	init   bool
	myType string

	Content  []string
	Dir      string
	FilePath string
}

func (r *TypeLicense) Init(dir string) *TypeLicense {
	r.init = true
	r.myType = "TypeLicense"
	prefix := r.myType + ".Init"

	r.Dir = dir
	r.FilePath = path.Join(r.Dir, Conf.FileLicense)
	helper.ReportDebug(r, prefix, false, true)
	if !helper.IsRegularFile(r.FilePath) {
		r.Err = errors.New(r.FilePath + " not found")
		helper.ErrsQueue(r.Err, prefix)
	}
	return r
}

func (r *TypeLicense) Dump() *TypeLicense {
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

func (r *TypeLicense) Update() *TypeLicense {
	prefix := r.myType + ".Update"
	if r.Err != nil {
		return r
	}
	if !r.init {
		r.Err = errors.New("not initialized")
		helper.ErrsQueue(r.Err, prefix)
	}
	if r.Err == nil {
		for lineNum := range r.Content {
			if r.Err == nil {
				helper.ReportDebug(&r.Content[lineNum], prefix, false, true)
				r.
					updateLicenseYear(&r.Content[lineNum])
				helper.ReportDebug(&r.Content[lineNum], prefix, false, true)
			}
		}
	}
	return r
}

// Read README.md into `Content`
func (r *TypeLicense) Read() *TypeLicense {
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
func (r *TypeLicense) Write() *TypeLicense {
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

// This work for:
//   - Copyright (c) xxxx
//   - Copyright © xxxx
func (r *TypeLicense) updateLicenseYear(line *string) *TypeLicense {
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
