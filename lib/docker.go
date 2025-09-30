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

	"github.com/J-Siu/go-basestruct"
	"github.com/J-Siu/go-helper/v2/cmd"
	"github.com/J-Siu/go-helper/v2/errs"
	"github.com/J-Siu/go-helper/v2/ezlog"
	"github.com/J-Siu/go-helper/v2/file"
	"github.com/J-Siu/go-helper/v2/str"
)

type TypeDocker struct {
	*basestruct.Base

	Content  []string
	Dir      string
	FilePath string
	Pkg      string

	Distro string
	Branch string
	Repo   []string

	VerCurr string
	VerNew  string

	Debug   bool
	Verbose bool
}

// Assuming branch = main + community
//
// Read and extract information from Dockerfile
func (t *TypeDocker) New(dir *string, debug, verbose bool) *TypeDocker {
	t.Base = new(basestruct.Base)
	t.Initialized = true
	t.MyType = "TypeDocker"
	prefix := t.MyType + ".Init"

	t.Verbose = verbose
	t.Repo = []string{"main", "community"}
	t.Dir = *dir
	t.FilePath = path.Join(t.Dir, "Dockerfile")
	ezlog.Debug().Nn(prefix).M(t).Out()
	if !file.IsRegularFile(t.FilePath) {
		t.Err = errors.New(t.FilePath + " not found")
	}
	if t.Err == nil {
		t.read().extract()
	}
	errs.Queue(prefix, t.Err)
	return t
}

func (t *TypeDocker) BuildTest() *TypeDocker {
	prefix := t.MyType + ".BuildTest"
	if t.Err != nil {
		return t
	}
	if !t.Initialized {
		t.Err = errors.New("not initialized")
		errs.Queue(prefix, t.Err)
	}
	var myCmd *cmd.Cmd
	imgName := t.Pkg + ":" + "auto_docker"
	if t.Err == nil {
		// RUN_CMD "docker build --quiet -t ${_img} ."
		args := []string{"build", "--quiet", "-t", imgName, "."}
		myCmd = cmd.Run("docker", &args, &t.Dir)
		t.Err = myCmd.Err
	}
	if t.Err == nil {
		// RUN_CMD "docker image rm ${_img}"
		args := []string{"image", "rm", imgName}
		myCmd = cmd.Run("docker", &args, &t.Dir)
		t.Err = myCmd.Err
	}
	if t.Verbose || t.Debug {
		if t.Err == nil {
			ezlog.Log().N(prefix).N(imgName).Msg("Success").Out()
		} else {
			ezlog.Log().N(prefix).N(imgName).Msg("Failed").Out()
		}
	}
	return t
}

func (t *TypeDocker) Dump() *TypeDocker {
	prefix := t.MyType + ".Dump"
	if t.Err != nil {
		return t
	}
	if !t.Initialized {
		t.Err = errors.New("not initialized")
		errs.Queue(prefix, t.Err)
	}
	if t.Err == nil {
		ezlog.Log().Nn(prefix).M(t).Out()
	}
	return t
}

// Update `Content`
func (t *TypeDocker) Update(dbAlpine *TypeDbAlpine) *TypeDocker {
	prefix := t.MyType + ".Update"
	if t.Err != nil {
		return t
	}
	if !t.Initialized {
		t.Err = errors.New("not initialized")
		errs.Queue(prefix, t.Err)
	}
	if t.Err == nil {
		// Check for new version
		for _, b := range t.Repo {
			verNew := *dbAlpine.PkgVerGet(t.Pkg, t.Branch, b)
			if dbAlpine.Err == nil {
				if verNew > t.VerNew {
					t.VerNew = verNew
					ezlog.Debug().N(prefix).N(t.Branch + "/" + b).N(t.Pkg).M(verNew).M(">").M(t.VerCurr).Out()
				}
			}
		}
		if t.VerNew > t.VerCurr {
			ezlog.Debug().N(prefix).N(t.Pkg).M(t.VerCurr).M("->").M(t.VerNew).Out()
			for index := range t.Content {
				t.Content[index] = strings.ReplaceAll(t.Content[index], t.VerCurr, t.VerNew)
			}
		}
	}
	return t
}

// Write `Content` to Dockerfile
func (t *TypeDocker) Write() *TypeDocker {
	prefix := t.MyType + ".Write"
	if t.Err != nil {
		return t
	}
	if !t.Initialized {
		t.Err = errors.New("not initialized")
		errs.Queue(prefix, t.Err)
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

// Read Dockerfile into `Content`
func (t *TypeDocker) read() *TypeDocker {
	prefix := t.MyType + ".read"
	if t.Err != nil {
		return t
	}
	if !t.Initialized {
		t.Err = errors.New("not initialized")
		errs.Queue(prefix, t.Err)
	}
	if t.Err == nil {
		t.Content, t.Err = file.ArrayRead(t.FilePath)
		if t.Err != nil {
			t.Err = errors.New(t.FilePath + " not found")
			errs.Queue(prefix, t.Err)
		}
	}
	return t
}

// Extract information from `Content`
func (t *TypeDocker) extract() *TypeDocker {
	prefix := t.MyType + ".extract"
	if t.Err != nil {
		return t
	}
	if !t.Initialized {
		t.Err = errors.New("not initialized")
		errs.Queue(prefix, t.Err)
	}
	if t.Err == nil {
		testing := "testing"
		branchTesting := t.Branch + "/" + testing
		for _, line := range t.Content {
			words := strings.Split(line, " ")
			switch strings.ToLower(words[0]) {
			case "from":
				ezlog.Debug().N(prefix).M(line).Out()
				t.Distro = words[1]
				// detect branch, eg. "alpine:edge" -> "edge"
				tmp := strings.Split(words[1], ":")
				if len(tmp) == 2 {
					t.Distro = tmp[0]
					t.Branch = tmp[1]
				}
			case "label":
				ezlog.Debug().N(prefix).M(words[1]).Out()
				label := strings.Split(words[1], "=")
				switch strings.ToLower(label[0]) {
				case "version":
					t.VerCurr = strings.ReplaceAll(label[1], "\"", "")
					t.VerNew = ""
				case "name":
					t.Pkg = strings.ReplaceAll(label[1], "\"", "")
				}
			default:
				// detect branch testing
				if strings.Contains(line, branchTesting) {
					if !str.ArrayContains(&t.Repo, &testing) {
						t.Repo = append(t.Repo, testing)
					}
				}
			}
		}
	}
	return t
}
