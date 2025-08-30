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
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"

	"github.com/J-Siu/go-helper"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const extTgz = ".tar.gz"

var DbAlpineDefault = TypeDbAlpine{
	Arch:   []string{"aarch64", "armhf", "armv7", "x86", "x86_64"},
	Branch: []string{"community", "main", "testing"},
	Distro: "alpine",
	Repo:   []string{"latest-stable", "edge"},
}

// Alpine package database struct base on repo, branch and arch
type TypeDbAlpine struct {
	Err    error
	init   bool
	myType string

	Db        *gorm.DB
	DirDb     string // full path of base database (db file + APKINDEX) directory
	FileDb    string // full path of the database file
	FileIndex string
	UrlBase   string

	Arch   []string
	Branch []string
	Distro string
	Repo   []string
}

type TypeDbAlpineRecord struct {
	// gorm.Model
	Repo   string `json:"Repo"`
	Branch string `json:"Branch"`
	Arch   string `json:"Arch"`
	Pkg    string `json:"Pkg"`
	Ver    string `json:"Ver"`
}

func (a *TypeDbAlpine) Init() *TypeDbAlpine {
	a.init = true
	a.myType = "TypeDbAlpine"
	prefix := a.myType + ".init"
	helper.ReportDebug("-- Start", prefix, false, true)

	a.Arch = DbAlpineDefault.Arch
	a.Branch = DbAlpineDefault.Branch
	a.Distro = DbAlpineDefault.Distro
	a.Repo = DbAlpineDefault.Repo

	a.DirDb = path.Join(Conf.DirCache, Conf.DirDB, a.Distro)
	a.FileDb = path.Join(Conf.DirCache, Conf.DirDB, a.Distro, a.Distro+".db")
	a.FileIndex = "APKINDEX"
	a.UrlBase = "http://dl-cdn.alpinelinux.org/alpine"

	helper.ReportDebug(a, prefix, false, false)
	helper.ReportDebug("-- End", prefix, false, true)

	return a
}

// DBConnect
//   - if database does not exist, an empty one will be created
func (a *TypeDbAlpine) DbConnect() *TypeDbAlpine {
	prefix := a.myType + ".DBConnect"
	helper.ReportDebug("-- Start", prefix, false, true)
	if a.Err != nil {
		return a
	}
	if !a.init {
		a.Err = errors.New("not initialized")
	}
	helper.ReportDebug(a.FileDb, prefix, false, true)

	if a.Err == nil {
		a.Err = os.MkdirAll(a.DirDb, os.ModePerm)
	}

	if a.Err == nil {
		a.Db, a.Err = gorm.Open(
			sqlite.Open(a.FileDb),
			&gorm.Config{
				QueryFields: true,
				Logger:      logger.Default.LogMode(logger.Silent),
			},
		)
		if a.Err != nil {
			a.Err = errors.New("cannot open " + a.FileDb)
		}
	}

	helper.ErrsQueue(a.Err, prefix)

	helper.ReportDebug("-- End", prefix, false, true)
	return a
}

// DbDump()
//   - This must be called after TypeDbAlpine.Init()
//   - DbDump DB to stdout
func (a *TypeDbAlpine) DbDump() *TypeDbAlpine {
	prefix := a.myType + ".DbDump"
	helper.ReportDebug("-- Start", prefix, false, true)

	if a.Err != nil {
		return a
	}
	if !a.init {
		a.Err = errors.New("not initialized")
	}

	if a.Db == nil {
		a.Err = errors.New("database not connected")
	}

	var rows helper.MyArray[TypeDbAlpineRecord]

	if a.Err == nil {
		result := a.Db.
			// Unscoped().
			Select([]string{"Pkg", "Ver", "Repo", "Branch", "Arch"}).
			Find(&rows)
		a.Err = result.Error
	}

	if a.Err == nil {
		helper.Report(rows, "", false, false)
		helper.Report(len(rows), "Rows", false, true)
	}

	helper.ErrsQueue(a.Err, prefix)

	helper.ReportDebug("-- End", prefix, false, true)
	return a
}

// Return immediately on error
func (a *TypeDbAlpine) DbUpdate() *TypeDbAlpine {
	prefix := a.myType + ".DbUpdate"
	helper.ReportDebug("-- Start", prefix, false, true)

	if a.Err != nil {
		return a
	}
	if !a.init {
		a.Err = errors.New("not initialized")
	}

	if a.Err == nil {
		a.Err = os.RemoveAll(a.DirDb) // Delete first
	}
	if a.Err == nil {
		a.DbConnect()
	}
	if a.Err == nil {
		a.Db.AutoMigrate(&TypeDbAlpineRecord{})
	}
	if a.Err == nil {
		a.Err = a.dbDownload()
	}

	helper.ErrsQueue(a.Err, prefix)

	helper.ReportDebug("-- End", prefix, false, true)
	return a
}

// PkgSearch
//   - Output result to stdout
//   - Return immediately on error
func (a *TypeDbAlpine) PkgSearch(pkg string) *TypeDbAlpine {
	prefix := a.myType + ".PkgSearch"
	helper.ReportDebug("-- Start", prefix, false, true)

	if a.Err != nil {
		return a
	}
	if !a.init {
		a.Err = errors.New("not initialized")
	}

	if a.Db == nil {
		a.Err = errors.New("database not connected")
	}

	var rows []TypeDbAlpineRecord

	if a.Err == nil {
		result := a.Db.
			Unscoped().
			Select([]string{"Pkg", "Ver", "Repo", "Branch", "Arch"})
		if FlagDbSearch.Exact {
			result = result.Where(map[string]interface{}{"Pkg": pkg})
		} else {
			result = result.Where("Pkg LIKE ?", "%"+pkg+"%")
		}
		result = result.Find(&rows)
		a.Err = result.Error
	}

	if a.Err == nil {
		for _, r := range rows {
			helper.Report(r.Pkg+" "+r.Ver+" "+r.Repo+" "+r.Branch+" "+r.Arch, "", true, false)
		}
	}

	helper.ErrsQueue(a.Err, prefix)

	helper.ReportDebug("-- End", prefix, false, true)
	return a
}

func (a *TypeDbAlpine) PkgVerGet(pkg string, repo string, branch string) (ver *string) {
	prefix := a.myType + ".PkgVerGet"
	helper.ReportDebug("-- Start", prefix, false, true)

	if a.Err != nil {
		return nil
	}
	if !a.init {
		a.Err = errors.New("not initialized")
	}

	if a.Db == nil {
		a.Err = errors.New("database not connected")
	}

	var row TypeDbAlpineRecord

	if a.Err == nil {
		result := a.Db.
			Unscoped().
			Where(map[string]interface{}{
				"Repo":   repo,
				"Branch": branch,
				"Arch":   a.Arch,
				"Pkg":    pkg,
			}).
			Find(&row)
		a.Err = result.Error
	}

	helper.ErrsQueue(a.Err, prefix)

	helper.ReportDebug("-- End", prefix, false, true)
	return &row.Ver
}

// Wrapper for Alpine APKINDEX download and database create/update
func (a *TypeDbAlpine) dbDownload() (err error) {
	prefix := a.myType + ".dbDownload"
	helper.ReportDebug("-- Start", prefix, false, true)

	for _, repo := range a.Repo {
		for _, branch := range a.Branch {
			for _, arch := range a.Arch {
				if !(repo == "latest-stable" && branch == "testing") {
					// Download APKINDEX.tar.gz
					err = a.idxDownload(repo, branch, arch)
					// Update database
					if err == nil {
						err = a.idx2db(repo, branch, arch)
					}
					helper.ErrsQueue(err, prefix)
				}
			}
		}
	}

	helper.ReportDebug("-- End", prefix, false, true)
	return err
}

func (a *TypeDbAlpine) idxDownload(repo string, branch string, arch string) (err error) {
	prefix := a.myType + ".idxDownload"
	helper.ReportDebug("-- Start", prefix, false, true)

	// Prepare download URL
	urlApkIndex, err := url.JoinPath(a.UrlBase, repo, branch, arch, a.FileIndex+extTgz)

	// Create directory
	dirArch := a.idxDir(repo, branch, arch)
	if err == nil {
		err = os.MkdirAll(dirArch, os.ModePerm)
		helper.ErrsQueue(err, prefix)
	}

	// Download APKINDEX.tar.gz
	filepathApkindex := path.Join(dirArch, a.FileIndex)
	filepathApkindexTgz := filepathApkindex + extTgz
	if err == nil {
		err = download(urlApkIndex, filepathApkindexTgz)
	}
	// Decompress APKINDEX.tar.gz
	if err == nil {
		err = untar(dirArch, filepathApkindexTgz)
	}

	helper.ReportDebug("-- End", prefix, false, true)
	return err
}

func (a *TypeDbAlpine) idx2db(repo string, branch string, arch string) (err error) {
	prefix := a.myType + ".idx2db"
	helper.ReportDebug("-- Start", prefix, false, true)

	dirArch := a.idxDir(repo, branch, arch)
	filepathApkIndex := path.Join(dirArch, a.FileIndex)
	helper.ReportDebug(filepathApkIndex, prefix, false, true)

	var rows helper.MyArray[TypeDbAlpineRecord]

	// Read APKINDEX file
	byteRead, err := os.ReadFile(filepathApkIndex)

	if err == nil {
		// Prepare DB rows
		lines := strings.Split(string(byteRead), "\n")
		var recordP *TypeDbAlpineRecord
		for _, l := range lines {
			if len(l) > 0 {
				switch l[0] {
				case 'P':
					recordP = &TypeDbAlpineRecord{
						Branch: branch,
						Repo:   repo,
						Arch:   arch,
						Pkg:    l[2:],
					}
				case 'V':
					recordP.Ver = l[2:]
					rows.Add(*recordP)
				}
			}
		}
		// Batch insert into DB
		result := a.Db.CreateInBatches(rows, 1000)
		err = result.Error
	}

	helper.ErrsQueue(err, prefix)
	helper.ReportDebug("-- End", prefix, false, true)
	return err
}

// Calculate(join) APKINDEX directory path base on `repo`, `branch`, `arch` and
func (a *TypeDbAlpine) idxDir(repo string, branch string, arch string) string {
	return path.Join(a.DirDb, repo, branch, arch)
}

// URL download to file
func download(url string, filepath string) (err error) {
	prefix := "download"
	helper.ReportDebug("-- Start", prefix, false, true)
	helper.ReportDebug("\t"+url+"->\n\t\t"+filepath, prefix, false, false)

	var res *http.Response
	out, err := os.Create(filepath)
	if err == nil {
		defer out.Close()
		res, err = http.Get(url)
		if res.Status[0:1] == "4" { // eg. 404
			err = errors.New(url + " " + res.Status)
		}
	}
	if err == nil {
		defer res.Body.Close()
		_, err = io.Copy(out, res.Body)
	}
	helper.ErrsQueue(err, prefix)

	helper.ReportDebug("-- End", prefix, false, true)
	return err
}

// Use command line tar to untar
func untar(dir string, filepath string) error {
	prefix := "untar"
	helper.ReportDebug("-- Start", prefix, false, true)

	opts := []string{"xf", filepath, "-C", dir}
	myCmd := helper.MyCmdRun("tar", &opts, nil)
	helper.ErrsQueue(myCmd.Err, prefix)

	helper.ReportDebug("-- End", prefix, false, true)
	return myCmd.Err
}
