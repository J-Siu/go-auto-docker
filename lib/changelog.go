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
	"strings"

	"github.com/J-Siu/go-helper/v2/basestruct"
	"github.com/J-Siu/go-helper/v2/errs"
	"github.com/J-Siu/go-helper/v2/ezlog"
	"github.com/J-Siu/go-helper/v2/file"
)

type TypeChangeLogProperty struct {
	Dir           *string `json:"Dir"`
	FileChangeLog *string `json:"FileChangeLog"` // CHANGELOG.md filename
	Pkg           *string `json:"Pkg"`
	VerCurr       *string `json:"VerCurr"`
	VerNew        *string `json:"VerNew"`
}

type TypeChangeLog struct {
	*basestruct.Base
	*TypeChangeLogProperty
	Content  *[]string `json:"Content"`
	FilePath string    `json:"FilePath"` //README path
}

func (t *TypeChangeLog) New(property *TypeChangeLogProperty) *TypeChangeLog {
	t.Base = new(basestruct.Base)
	t.TypeChangeLogProperty = property
	t.MyType = "TypeChangeLog"

	t.FilePath = path.Join(*t.Dir, *t.FileChangeLog)

	prefix := t.MyType + ".New"
	ezlog.Debug().N(prefix).Lm(t).Out()
	if !file.IsRegularFile(t.FilePath) {
		t.Err = errors.New(t.FilePath + " not found")
		errs.Queue(prefix, t.Err)
	}
	t.Initialized = true
	return t
}

func (t *TypeChangeLog) Dump() *TypeChangeLog {
	prefix := t.MyType + ".Dump"
	if !t.CheckErrInit(prefix) {
		errs.Queue(prefix, t.Err)
	}
	if t.Err == nil {
		ezlog.Log().N(prefix).Lm(t).Out()
	}
	return t
}

func (t *TypeChangeLog) Update() *TypeChangeLog {
	prefix := t.MyType + ".Update"
	var contentNew []string
	if !t.CheckErrInit(prefix) {
		return t
	}
	if *t.VerNew > *t.VerCurr {
		ezlog.Debug().N(prefix).N(t.Pkg).M(t.VerCurr).M("->").M(t.VerNew).Out()
		for _, line := range *t.Content {
			ezlog.Debug().N(prefix).M(&line).Out()
			if strings.Contains(line, *t.VerNew) {
				t.Err = errors.New(*t.FileChangeLog + " contains " + *t.VerNew)
				break
			} else if line != "" {
				contentNew = append(contentNew, line)
			}
		}
	} else {
		t.Err = errors.New(": Version not newer")
	}
	if t.Err == nil {
		contentNew = append(contentNew, "- "+*t.VerNew)
		contentNew = append(contentNew, "  - Auto update to "+*t.VerNew)
		t.Content = &contentNew
	}
	errs.Queue(prefix, t.Err)
	return t
}

// Read README.md into `Content`
func (t *TypeChangeLog) Read() *TypeChangeLog {
	prefix := t.MyType + ".Read"
	if t.Err != nil {
		return t
	}
	if !t.Initialized {
		t.Err = errors.New("not initialized")
	}
	if t.Err == nil {
		t.Content, t.Err = file.ReadStrArray(t.FilePath)
		if t.Err != nil {
			t.Err = errors.New(t.FilePath + " not found")
		} else {
		}
	}
	errs.Queue(prefix, t.Err)
	return t
}

// Write `Content` into README.md
func (t *TypeChangeLog) Write() *TypeChangeLog {
	prefix := t.MyType + ".Write"
	if !t.CheckErrInit(prefix) {
		return t
	}
	if t.Err == nil {
		fileStats, err := os.Stat(t.FilePath)
		if err != nil {
			// Should never happen at this stage, but ...
			t.Err = err
		} else {
			file.WriteStrArray(t.FilePath, t.Content, fileStats.Mode())
		}
	}
	errs.Queue(prefix, t.Err)
	return t
}
