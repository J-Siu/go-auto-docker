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

	"github.com/J-Siu/go-basestruct"
	"github.com/J-Siu/go-helper/v2/errs"
	"github.com/J-Siu/go-helper/v2/ezlog"
	"github.com/J-Siu/go-helper/v2/file"
	"github.com/J-Siu/go-helper/v2/strany"
)

type TypeReadme struct {
	*basestruct.Base

	Content  []string `json:"content,omitempty"`
	Dir      string   `json:"dir,omitempty"`
	FilePath string   `json:"file_path,omitempty"`

	TagReadmeLogStart string `json:"tag_readme_log_start,omitempty"`
	TagReadmeLogEnd   string `json:"tag_readme_log_end,omitempty"`

	Pkg     string `json:"pkg,omitempty"`
	VerCurr string `json:"ver_curr,omitempty"`
	VerNew  string `json:"ver_new,omitempty"`

	changeLogStart bool
	changeLogEnd   bool
	newLogHeader   string
	newLogContent  string
}

func (t *TypeReadme) New(dir, pkg, verCurr, verNew, fileReadme, tagReadmeLogStart, tagReadmeLogEnd *string) *TypeReadme {
	t.Base = new(basestruct.Base)
	t.Initialized = true
	t.MyType = "TypeReadme"
	t.Dir = *dir
	t.FilePath = path.Join(t.Dir, *fileReadme)
	t.TagReadmeLogStart = *tagReadmeLogStart
	t.TagReadmeLogEnd = *tagReadmeLogEnd

	t.Pkg = *pkg
	t.VerCurr = *verCurr
	t.VerNew = *verNew

	prefix := t.MyType + ".New"
	ezlog.Debug().Nn(prefix).M(t).Out()
	if !file.IsRegularFile(t.FilePath) {
		t.Err = errors.New("TypeReadme.Init: " + t.FilePath + " not found")
		errs.Queue(prefix, t.Err)
	}
	return t
}

func (t *TypeReadme) Dump() *TypeReadme {
	prefix := t.MyType + ".Dump"
	if !t.CheckErrInit(prefix) {
		errs.Queue(prefix, t.Err)
	}
	if t.Err == nil {
		ezlog.Log().Nn(prefix).M(t).Out()
	}
	return t
}

func (t *TypeReadme) Update() *TypeReadme {
	prefix := t.MyType + ".Update"
	if !t.CheckErrInit(prefix) {
		errs.Queue(prefix, t.Err)
	}
	if t.Err == nil {
		if t.VerNew > t.VerCurr {
			ezlog.Debug().N(prefix).N(t.Pkg).M(t.VerCurr).M("->").M(t.VerNew).Out()
			t.changeLogStart = false
			t.changeLogEnd = false
			t.newLogHeader = "- " + t.VerNew
			t.newLogContent = "  - Auto update to " + t.VerNew
			// for lineNum := range r.Content {
			for lineNum := range t.Content {
				if t.Err == nil {
					ezlog.Debug().N(prefix).M(&t.Content[lineNum]).Out()
					t.
						updateLog(&t.Content[lineNum]).
						updateLicenseYear(&t.Content[lineNum])
					ezlog.Debug().N(prefix).M(&t.Content[lineNum]).Out()
				}
			}
		}
	}
	return t
}

// Read README.md into `Content`
func (t *TypeReadme) Read() *TypeReadme {
	prefix := t.MyType + ".Read"
	if t.Err != nil {
		return t
	}
	if !t.Initialized {
		t.Err = errors.New("not initialized")
	}
	if t.Err == nil {
		t.Content, t.Err = file.ArrayRead(t.FilePath)
		if t.Err != nil {
			t.Err = errors.New(t.FilePath + " not found")
		} else {
		}
	}
	errs.Queue(prefix, t.Err)
	return t
}

// Write `Content` into README.md
func (t *TypeReadme) Write() *TypeReadme {
	prefix := t.MyType + ".Write"
	if t.Err != nil {
		return t
	}
	if !t.Initialized {
		t.Err = errors.New("not initialized")
	}
	if t.Err == nil {
		fileStats, err := os.Stat(t.FilePath)
		if err != nil {
			// Should never happen at this stage, but ...
			t.Err = err
		} else {
			file.ArrayWrite(t.FilePath, t.Content, fileStats.Mode())
		}
	}
	errs.Queue(prefix, t.Err)
	return t
}

func (t *TypeReadme) updateLog(line *string) *TypeReadme {
	prefix := t.MyType + ".updateLog"
	if strings.EqualFold(*line, t.TagReadmeLogStart) {
		t.changeLogStart = true
		ezlog.Debug().N(prefix).TxtStart().Out()
	}
	if !t.changeLogEnd && t.changeLogStart {
		ezlog.Debug().N(prefix).M("Change").TxtStart().M(line).Out()
		if strings.Contains(*line, t.newLogHeader) {
			// Something wrong if new log header already exist
			t.Err = errors.New("change log already contains " + t.VerNew)
			errs.Queue(prefix, t.Err)
		}
		if t.Err == nil && strings.Contains(*line, t.TagReadmeLogEnd) {
			*line = t.newLogHeader + "\n" + t.newLogContent + "\n" + *line
			ezlog.Debug().N(prefix).M(line).Out()
		}
	}
	if strings.Contains(*line, t.TagReadmeLogEnd) {
		t.changeLogEnd = true
		ezlog.Debug().N(prefix).M("Change").TxtEnd().Out()
	}
	ezlog.Debug().N(prefix).M(line).Out()
	return t
}

// This work for:
//   - Copyright (c) xxxx
//   - Copyright © xxxx
func (t *TypeReadme) updateLicenseYear(line *string) *TypeReadme {
	prefix := t.MyType + ".updateLicenseYear"
	c := []string{"(c)", "©"}
	cRaw := []string{`\(c\)`, `©`}
	copyright := "Copyright"
	year := *strany.Any(time.Now().Year())
	for i := range c {
		// "(?i)" <- means case insensitive
		re := regexp.MustCompile("(?i)" + copyright + " " + cRaw[i] + ` \d\d\d\d`)
		*line = re.ReplaceAllString(*line, copyright+" "+c[i]+" "+year)
	}
	ezlog.Debug().N(prefix).M(line).Out()
	return t
}
