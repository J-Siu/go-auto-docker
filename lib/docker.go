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

	"github.com/J-Siu/go-helper"
)

type TypeDocker struct {
	Err    error
	init   bool
	myType string

	Content  []string
	Dir      string
	FilePath string
	Pkg      string

	Distro string
	Branch string
	Repo   []string

	VerCurr string
	VerNew  string
}

// Assuming branch = main + community
//
// Read and extract information from Dockerfile
func (d *TypeDocker) Init(dir string) *TypeDocker {
	d.init = true
	d.myType = "TypeDocker"
	prefix := d.myType + ".Init"

	d.Repo = []string{"main", "community"}
	d.Dir = dir
	d.FilePath = path.Join(d.Dir, "Dockerfile")
	helper.ReportDebug(d, prefix, false, true)
	if !helper.IsRegularFile(d.FilePath) {
		d.Err = errors.New(d.FilePath + " not found")
	}
	if d.Err == nil {
		d.read().extract()
	}
	helper.ErrsQueue(d.Err, prefix)
	return d
}

func (d *TypeDocker) BuildTest() *TypeDocker {
	prefix := d.myType + ".BuildTest"
	if d.Err != nil {
		return d
	}
	if !d.init {
		d.Err = errors.New("not initialized")
		helper.ErrsQueue(d.Err, prefix)
	}
	var cmd *helper.MyCmd
	imgName := d.Pkg + ":" + "auto_docker"
	if d.Err == nil {
		// RUN_CMD "docker build --quiet -t ${_img} ."
		args := []string{"build", "--quiet", "-t", imgName, "."}
		cmd = helper.MyCmdRun("docker", &args, &d.Dir)
		d.Err = cmd.Err
	}
	if d.Err == nil {
		// RUN_CMD "docker image rm ${_img}"
		args := []string{"image", "rm", imgName}
		cmd = helper.MyCmdRun("docker", &args, &d.Dir)
		d.Err = cmd.Err
	}
	if Flag.Verbose || Flag.Debug {
		if d.Err == nil {
			helper.Report(imgName+": Success", prefix, false, true)
		} else {
			helper.Report(imgName+": Failed", prefix, false, true)
		}
	}
	return d
}

func (d *TypeDocker) Dump() *TypeDocker {
	prefix := d.myType + ".Dump"
	if d.Err != nil {
		return d
	}
	if !d.init {
		d.Err = errors.New("not initialized")
		helper.ErrsQueue(d.Err, prefix)
	}
	if d.Err == nil {
		helper.Report(d, prefix, false, false)
	}
	return d
}

// Update `Content`
func (d *TypeDocker) Update() *TypeDocker {
	prefix := d.myType + ".Update"
	if d.Err != nil {
		return d
	}
	if !d.init {
		d.Err = errors.New("not initialized")
		helper.ErrsQueue(d.Err, prefix)
	}
	if d.Err == nil {
		// Check for new version
		for _, b := range d.Repo {
			verNew := *DbAlpine.PkgVerGet(d.Pkg, d.Branch, b)
			if DbAlpine.Err == nil {
				if verNew > d.VerNew {
					d.VerNew = verNew
					helper.ReportDebug(d.Branch+"/"+b+":"+d.Pkg+":"+verNew+">"+d.VerCurr, prefix, false, true)
				}
			}
		}
		if d.VerNew > d.VerCurr {
			helper.ReportDebug(d.Pkg+": "+d.VerCurr+" -> "+d.VerNew, prefix, false, true)
			for index := range d.Content {
				d.Content[index] = strings.ReplaceAll(d.Content[index], d.VerCurr, d.VerNew)
			}
		}
	}
	return d
}

// Write `Content` to Dockerfile
func (d *TypeDocker) Write() *TypeDocker {
	prefix := d.myType + ".Write"
	if d.Err != nil {
		return d
	}
	if !d.init {
		d.Err = errors.New("not initialized")
		helper.ErrsQueue(d.Err, prefix)
	}
	if d.Err == nil {
		fileStats, err := os.Stat(d.FilePath)
		if err != nil {
			// Should never happen at this stage, but ...
			d.Err = err
		} else {
			helper.FileStrArrWrite(d.FilePath, d.Content, fileStats.Mode())
		}
	}
	helper.ErrsQueue(d.Err, prefix)
	return d
}

// Read Dockerfile into `Content`
func (d *TypeDocker) read() *TypeDocker {
	prefix := d.myType + ".read"
	if d.Err != nil {
		return d
	}
	if !d.init {
		d.Err = errors.New("not initialized")
		helper.ErrsQueue(d.Err, prefix)
	}
	if d.Err == nil {
		d.Content, d.Err = helper.FileStrArrRead(d.FilePath)
		if d.Err != nil {
			d.Err = errors.New(d.FilePath + " not found")
			helper.ErrsQueue(d.Err, prefix)
		}
	}
	return d
}

// Extract information from `Content`
func (d *TypeDocker) extract() *TypeDocker {
	prefix := d.myType + ".extract"
	if d.Err != nil {
		return d
	}
	if !d.init {
		d.Err = errors.New("not initialized")
		helper.ErrsQueue(d.Err, prefix)
	}
	if d.Err == nil {
		testing := "testing"
		branchTesting := d.Branch + "/" + testing
		for _, line := range d.Content {
			words := strings.Split(line, " ")
			switch strings.ToLower(words[0]) {
			case "from":
				helper.ReportDebug(line, prefix, false, true)
				d.Distro = words[1]
				// detect branch, eg. "alpine:edge" -> "edge"
				tmp := strings.Split(words[1], ":")
				if len(tmp) == 2 {
					d.Distro = tmp[0]
					d.Branch = tmp[1]
				}
			case "label":
				helper.ReportDebug(words[1], prefix, false, true)
				label := strings.Split(words[1], "=")
				switch strings.ToLower(label[0]) {
				case "version":
					d.VerCurr = strings.ReplaceAll(label[1], "\"", "")
					d.VerNew = ""
				case "name":
					d.Pkg = strings.ReplaceAll(label[1], "\"", "")
				}
			default:
				// detect branch testing
				if strings.Contains(line, branchTesting) {
					if !helper.StrArrayPtrContain(&d.Repo, &testing) {
						d.Repo = append(d.Repo, testing)
					}
				}
			}
		}
	}
	return d
}
