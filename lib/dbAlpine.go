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

	"github.com/J-Siu/go-helper/v2/array"
	"github.com/J-Siu/go-helper/v2/basestruct"
	"github.com/J-Siu/go-helper/v2/cmd"
	"github.com/J-Siu/go-helper/v2/errs"
	"github.com/J-Siu/go-helper/v2/ezlog"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const extTgz = ".tar.gz"

var DbAlpineDefault = TypeDbAlpine{
	FileIndex: "APKINDEX",
	UrlBase:   "http://dl-cdn.alpinelinux.org/alpine",

	Distro:     "alpine",
	Branch:     []string{"latest-stable", "edge"},
	Repository: []string{"community", "main", "testing"},
	Arch:       []string{"aarch64", "armhf", "armv7", "x86", "x86_64"},
}

// Alpine package database struct base on repo, branch and arch
type TypeDbAlpine struct {
	*basestruct.Base

	Db        *gorm.DB
	DirDb     string // full path of base database (db file + APKINDEX) directory
	FileDb    string // full path of the database file
	FileIndex string
	UrlBase   string

	Distro     string
	Branch     []string
	Repository []string
	Arch       []string
}

type TypeDbAlpineRecord struct {
	// gorm.Model
	Pkg    string `json:"Pkg"`
	Branch string `json:"Branch"`
	Repo   string `json:"Repo"`
	Arch   string `json:"Arch"`
	Ver    string `json:"Ver"`
}

func (t *TypeDbAlpine) New(dirCache, dirDb *string, alpineBranch *[]string) *TypeDbAlpine {
	t.Base = new(basestruct.Base)
	t.Initialized = true
	t.MyType = "TypeDbAlpine"
	prefix := t.MyType + ".init"
	ezlog.Debug().N(prefix).TxtStart().Out()

	t.DirDb = path.Join(*dirCache, *dirDb, t.Distro)
	t.FileDb = path.Join(*dirCache, *dirDb, t.Distro, t.Distro+".db")

	t.setDefault(alpineBranch)

	ezlog.Debug().Nn(prefix).M(t).Out()

	ezlog.Debug().N(prefix).TxtEnd().Out()
	return t
}

func (t *TypeDbAlpine) setDefault(alpineBranch *[]string) *TypeDbAlpine {
	t.Arch = DbAlpineDefault.Arch
	t.Repository = DbAlpineDefault.Repository
	t.Distro = DbAlpineDefault.Distro
	t.Branch = DbAlpineDefault.Branch
	t.FileIndex = DbAlpineDefault.FileIndex
	t.UrlBase = DbAlpineDefault.UrlBase

	if len(*alpineBranch) == 0 {
		t.Branch = DbAlpineDefault.Branch
	} else {
		t.Branch = *alpineBranch
	}

	return t
}

// DBConnect
//   - if database does not exist, an empty one will be created
func (t *TypeDbAlpine) DbConnect() *TypeDbAlpine {
	prefix := t.MyType + ".DBConnect"
	ezlog.Debug().N(prefix).TxtStart().Out()
	if t.Err != nil {
		return t
	}
	if !t.Initialized {
		t.Err = errors.New("not initialized")
	}
	ezlog.Debug().N(prefix).M(t.FileDb).Out()

	if t.Err == nil {
		t.Err = os.MkdirAll(t.DirDb, os.ModePerm)
	}

	if t.Err == nil {
		t.Db, t.Err = gorm.Open(
			sqlite.Open(t.FileDb),
			&gorm.Config{
				QueryFields: true,
				Logger:      logger.Default.LogMode(logger.Silent),
			},
		)
		if t.Err != nil {
			t.Err = errors.New("cannot open " + t.FileDb)
		}
	}

	errs.Queue(prefix, t.Err)

	ezlog.Debug().N(prefix).TxtEnd().Out()
	return t
}

// DbDump()
//   - This must be called after TypeDbAlpine.New()
//   - DbDump DB to stdout
func (t *TypeDbAlpine) DbDump() *TypeDbAlpine {
	prefix := t.MyType + ".DbDump"
	ezlog.Debug().N(prefix).TxtStart().Out()

	if t.Err != nil {
		return t
	}
	if !t.Initialized {
		t.Err = errors.New("not initialized")
	}

	if t.Db == nil {
		t.Err = errors.New("database not connected")
	}

	var rows array.Array[TypeDbAlpineRecord]

	if t.Err == nil {
		result := t.Db.
			// Unscoped().
			Select([]string{"Pkg", "Ver", "Repo", "Branch", "Arch"}).
			Find(&rows)
		t.Err = result.Error
	}

	if t.Err == nil {
		ezlog.Log().M(rows).Out()
		ezlog.Log().N("Rows").M(len(rows)).Out()
	}

	errs.Queue(prefix, t.Err)

	ezlog.Debug().N(prefix).TxtEnd().Out()
	return t
}

// Return immediately on error
func (t *TypeDbAlpine) DbUpdate() *TypeDbAlpine {
	prefix := t.MyType + ".DbUpdate"
	ezlog.Debug().N(prefix).TxtStart().Out()

	if t.Err != nil {
		return t
	}
	if !t.Initialized {
		t.Err = errors.New("not initialized")
	}

	if t.Err == nil {
		t.Err = os.RemoveAll(t.DirDb) // Delete first
	}
	if t.Err == nil {
		t.DbConnect()
	}
	if t.Err == nil {
		t.Db.AutoMigrate(&TypeDbAlpineRecord{})
	}
	if t.Err == nil {
		t.Err = t.dbDownload()
	}

	errs.Queue(prefix, t.Err)

	ezlog.Debug().N(prefix).TxtEnd().Out()
	return t
}

// PkgSearch
//   - Output result to stdout
//   - Return immediately on error
func (t *TypeDbAlpine) PkgSearch(pkg string, exact bool) *TypeDbAlpine {
	prefix := t.MyType + ".PkgSearch"
	ezlog.Debug().N(prefix).TxtStart().Out()

	if t.Err != nil {
		return t
	}
	if !t.Initialized {
		t.Err = errors.New("not initialized")
	}

	if t.Db == nil {
		t.Err = errors.New("database not connected")
	}

	var rows []TypeDbAlpineRecord

	if t.Err == nil {
		result := t.Db.
			Unscoped().
			Select([]string{"Pkg", "Ver", "Branch", "Repo", "Arch"})
		if exact {
			result = result.Where(map[string]interface{}{"Pkg": pkg})
		} else {
			result = result.Where("Pkg LIKE ?", "%"+pkg+"%")
		}
		result = result.Find(&rows)
		t.Err = result.Error
	}

	if t.Err == nil {
		for _, r := range rows {
			ezlog.Log().M(r.Pkg).M(r.Ver).M(r.Repo).M(r.Branch).M(r.Arch).Out()
		}
	}

	errs.Queue(prefix, t.Err)

	ezlog.Debug().N(prefix).TxtEnd().Out()
	return t
}

func (t *TypeDbAlpine) PkgVerGet(pkg string, branch string, repo string) (ver *string) {
	prefix := t.MyType + ".PkgVerGet"
	ezlog.Debug().N(prefix).TxtStart().Out()

	if t.Err != nil {
		return nil
	}
	if !t.Initialized {
		t.Err = errors.New("not initialized")
	}

	if t.Db == nil {
		t.Err = errors.New("database not connected")
	}

	var row TypeDbAlpineRecord

	if t.Err == nil {
		result := t.Db.
			Unscoped().
			Where(map[string]interface{}{
				"Branch": branch,
				"Repo":   repo,
				"Arch":   t.Arch,
				"Pkg":    pkg,
			}).
			Find(&row)
		t.Err = result.Error
	}

	errs.Queue(prefix, t.Err)

	ezlog.Debug().N(prefix).TxtEnd().Out()
	return &row.Ver
}

// Wrapper for Alpine APKINDEX download and database create/update
func (t *TypeDbAlpine) dbDownload() (err error) {
	prefix := t.MyType + ".dbDownload"
	ezlog.Debug().N(prefix).TxtStart().Out()

	for _, branch := range t.Branch {
		for _, repo := range t.Repository {
			for _, arch := range t.Arch {
				// stable branches don't have "testing"
				stable := branch == "latest-stable" || strings.ToLower(branch)[0] == 'v'
				if !(stable && repo == "testing") {
					// Download APKINDEX.tar.gz
					err = t.idxDownload(branch, repo, arch)
					// Update database
					if err == nil {
						err = t.idx2db(branch, repo, arch)
					}
					errs.Queue(prefix, err)
				}
			}
		}
	}

	ezlog.Debug().N(prefix).TxtEnd().Out()
	return err
}

func (t *TypeDbAlpine) idxDownload(branch string, repo string, arch string) (err error) {
	prefix := t.MyType + ".idxDownload"
	ezlog.Debug().N(prefix).TxtStart().Out()

	// Prepare download URL
	urlApkIndex, err := url.JoinPath(t.UrlBase, branch, repo, arch, t.FileIndex+extTgz)

	// Create directory
	dirArch := t.idxDir(branch, repo, arch)
	if err == nil {
		err = os.MkdirAll(dirArch, os.ModePerm)
		errs.Queue(prefix, err)
	}

	// Download APKINDEX.tar.gz
	filepathApkindex := path.Join(dirArch, t.FileIndex)
	filepathApkindexTgz := filepathApkindex + extTgz
	if err == nil {
		err = download(urlApkIndex, filepathApkindexTgz)
	}
	// Decompress APKINDEX.tar.gz
	if err == nil {
		err = untar(dirArch, filepathApkindexTgz)
	}

	ezlog.Debug().N(prefix).TxtEnd().Out()
	return err
}

func (t *TypeDbAlpine) idx2db(branch string, repo string, arch string) (err error) {
	prefix := t.MyType + ".idx2db"
	ezlog.Debug().N(prefix).TxtStart().Out()

	dirArch := t.idxDir(branch, repo, arch)
	filepathApkIndex := path.Join(dirArch, t.FileIndex)
	ezlog.Debug().N(prefix).M(filepathApkIndex).Out()

	var rows array.Array[TypeDbAlpineRecord]

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
		result := t.Db.CreateInBatches(rows, 1000)
		err = result.Error
	}

	errs.Queue(prefix, err)
	ezlog.Debug().N(prefix).TxtEnd().Out()
	return err
}

// Calculate(join) APKINDEX directory path base on `repo`, `branch`, `arch` and
func (t *TypeDbAlpine) idxDir(branch string, repo string, arch string) string {
	return path.Join(t.DirDb, branch, repo, arch)
}

// URL download to file
func download(url string, filepath string) (err error) {
	prefix := "download"
	ezlog.Debug().N(prefix).TxtStart().Out()
	ezlog.Debug().Nn(prefix).T().Mn(url).T().T().M(filepath)

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
	// errs.Queue(prefix, err)

	ezlog.Debug().N(prefix).TxtEnd().Out()
	return err
}

// Use command line tar to untar
func untar(dir string, filepath string) error {
	prefix := "untar"
	ezlog.Debug().N(prefix).TxtStart().Out()

	opts := []string{"xf", filepath, "-C", dir}
	myCmd := cmd.Run("tar", &opts, nil)
	errs.Queue(prefix, myCmd.Err)

	ezlog.Debug().N(prefix).TxtEnd().Out()
	return myCmd.Err
}
