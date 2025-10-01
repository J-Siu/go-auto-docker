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

	"github.com/J-Siu/go-helper/v2/basestruct"
	"github.com/J-Siu/go-helper/v2/errs"
	"github.com/J-Siu/go-helper/v2/ezlog"
	"github.com/J-Siu/go-helper/v2/file"
	"github.com/go-git/go-git/v6"
	"github.com/go-git/go-git/v6/config"
	"github.com/go-git/go-git/v6/plumbing"
)

type TypeRepository struct {
	*basestruct.Base

	DirCache     string // project copy location
	DirCacheBase string
	DirSrc       string // project(git repo) original location
	Name         string

	Verbose bool
}

func (t *TypeRepository) New(workPath, dirCache, dirRepo *string, verbose bool) *TypeRepository {
	t.Base = new(basestruct.Base)
	t.MyType = "TypeRepository"
	prefix := t.MyType + ".init"
	ezlog.Debug().N(prefix).TxtStart().Out()

	t.Verbose = verbose
	t.DirSrc = *workPath
	if t.DirSrc == "." {
		t.DirSrc = *file.CurrentPath()
	}

	_, t.Name = path.Split(t.DirSrc)
	t.DirCacheBase = path.Join(*dirCache, *dirRepo)
	t.DirCache = path.Join(t.DirCacheBase, t.Name)

	ezlog.Debug().N(prefix).M(t).Out()

	_, t.Err = git.PlainOpen(t.DirSrc)
	if t.Err != nil {
		t.Err = errors.New(t.DirSrc + " is not a git repository.")
		errs.Queue(prefix, t.Err)
	}

	t.Initialized = true

	return t
}

func (t *TypeRepository) CopySrcToCache() *TypeRepository {
	prefix := t.MyType + ".CopySrcToCache"
	ezlog.Debug().N(prefix).TxtStart().Out()
	if t.Err != nil {
		return t
	}
	if !t.Initialized {
		t.Err = errors.New("not initialized")
	}
	if t.Err == nil {
		t.copyDir(t.DirSrc, t.DirCache)
	}
	return t
}

func (t *TypeRepository) CopyCacheToSrc() *TypeRepository {
	prefix := t.MyType + ".CopyCacheToSrc"
	ezlog.Debug().N(prefix).TxtStart().Out()
	if t.Err != nil {
		return t
	}
	if !t.Initialized {
		t.Err = errors.New("not initialized")
	}
	if t.Err == nil {
		t.copyDir(t.DirCache, t.DirSrc)
	}
	return t
}

func (t *TypeRepository) Commit(msg string, tag bool, cache bool) *TypeRepository {
	prefix := t.MyType + ".Commit"
	ezlog.Debug().N(prefix).TxtStart().Out()
	if t.Err != nil {
		return t
	}
	if !t.Initialized {
		t.Err = errors.New("not initialized")
	}

	commitOptions := git.CommitOptions{}
	// var commitObj *object.Commit
	var commit plumbing.Hash
	var gitConf *config.Config
	var gitDir string
	var gitRepo *git.Repository
	var gitWorktree *git.Worktree
	var gitHead *plumbing.Reference
	if cache {
		gitDir = t.DirCache
	} else {
		gitDir = t.DirSrc
	}
	prefix = prefix + "(" + gitDir + ")"
	// Repository open
	if t.Err == nil {
		gitRepo, t.Err = git.PlainOpen(gitDir)
	}
	// Repository load config
	if t.Err == nil {
		gitConf, t.Err = gitRepo.Config()
		gitConf.User.Name = "J" // TODO: why J?
		ezlog.Debug().Nn(prefix).M(gitConf).Out()
	}
	// Repository worktree
	if t.Err == nil {
		gitWorktree, t.Err = gitRepo.Worktree()
		ezlog.Debug().N(prefix).M("repo worktree created").Out()
	}
	// Worktree stage all
	if t.Err == nil {
		addOpt := git.AddOptions{
			All:  true,
			Path: gitDir,
		}
		t.Err = gitWorktree.AddWithOptions(&addOpt)
		ezlog.Debug().N(prefix).M("worktree staged(" + gitDir + ")").Out()
	}
	// Worktree commit
	if t.Err == nil {
		commit, t.Err = gitWorktree.Commit(msg, &commitOptions)
		// commit, p.Err = gitWorktree.Commit(msg, nil)
		ezlog.Debug().N(prefix).M("worktree committed").Out()
	}
	// Repository commit
	if t.Err == nil {
		// commitObj, p.Err = p.gitRepo.CommitObject(commit)
		_, t.Err = gitRepo.CommitObject(commit)
		ezlog.Debug().N(prefix).M("repo committed").Out()
	}
	// Repository tag
	if tag {
		if t.Err == nil {
			gitHead, t.Err = gitRepo.Head()
		}
		if t.Err == nil {
			_, t.Err = gitRepo.CreateTag(msg, gitHead.Hash(), nil)
			ezlog.Debug().N(prefix).M("tag(" + msg + ")").Out()
		}
	}
	errs.Queue(prefix, t.Err)
	return t
}

func (t *TypeRepository) copyDir(dirSrc, dirDest string) *TypeRepository {
	prefix := t.MyType + ".copyDir"
	ezlog.Debug().N(prefix).TxtStart().Out()
	if t.Err != nil {
		return t
	}
	if t.Err == nil {
		os.RemoveAll(dirDest) // Delete first
		t.Err = os.CopyFS(dirDest, os.DirFS(dirSrc))
	}
	if t.Err == nil && t.Verbose {
		ezlog.Log().N(prefix).M(dirSrc).M("->").M(dirDest).Out()
	}
	errs.Queue(prefix, t.Err)
	return t
}
