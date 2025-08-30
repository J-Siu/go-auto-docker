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

	"github.com/J-Siu/go-helper"
	"github.com/go-git/go-git/v6"
	"github.com/go-git/go-git/v6/config"
	"github.com/go-git/go-git/v6/plumbing"
)

type TypeRepository struct {
	Err    error
	init   bool
	myType string

	DirCache     string // project copy location
	DirCacheBase string
	DirSrc       string // project(git repo) original location
	Name         string
}

func (p *TypeRepository) Init(workPath string) *TypeRepository {
	p.myType = "TypeRepository"
	prefix := p.myType + ".init"
	helper.ReportDebug("-- Start", prefix, false, true)

	p.DirSrc = workPath
	if p.DirSrc == "." {
		p.DirSrc = *helper.CurrentPath()
	}

	_, p.Name = path.Split(p.DirSrc)
	p.DirCacheBase = path.Join(Conf.DirCache, Default.DirProj)
	p.DirCache = path.Join(p.DirCacheBase, p.Name)

	helper.ReportDebug(p, prefix, false, true)

	_, p.Err = git.PlainOpen(p.DirSrc)
	if p.Err != nil {
		p.Err = errors.New(p.DirSrc + " is not a git repository.")
		helper.ErrsQueue(p.Err, prefix)
	}

	p.init = true

	return p
}

func (p *TypeRepository) CopySrcToCache() *TypeRepository {
	prefix := p.myType + ".CopySrcToCache"
	helper.ReportDebug("-- Start", prefix, false, true)
	if p.Err != nil {
		return p
	}
	if !p.init {
		p.Err = errors.New("not initialized")
	}
	if p.Err == nil {
		p.copyDir(p.DirSrc, p.DirCache)
	}
	return p
}

func (p *TypeRepository) CopyCacheToSrc() *TypeRepository {
	prefix := p.myType + ".CopyCacheToSrc"
	helper.ReportDebug("-- Start", prefix, false, true)
	if p.Err != nil {
		return p
	}
	if !p.init {
		p.Err = errors.New("not initialized")
	}
	if p.Err == nil {
		p.copyDir(p.DirCache, p.DirSrc)
	}
	return p
}

func (p *TypeRepository) Commit(msg string, tag bool, cache bool) *TypeRepository {
	prefix := p.myType + ".Commit"
	helper.ReportDebug("-- Start", prefix, false, true)
	if p.Err != nil {
		return p
	}
	if !p.init {
		p.Err = errors.New("not initialized")
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
		gitDir = p.DirCache
	} else {
		gitDir = p.DirSrc
	}
	prefix = prefix + "(" + gitDir + ")"
	// Repository open
	if p.Err == nil {
		gitRepo, p.Err = git.PlainOpen(gitDir)
	}
	// Repository load config
	if p.Err == nil {
		gitConf, p.Err = gitRepo.Config()
		gitConf.User.Name = "J"
		helper.ReportDebug(gitConf, prefix, false, true)
	}
	// Repository worktree
	if p.Err == nil {
		gitWorktree, p.Err = gitRepo.Worktree()
		helper.ReportDebug("repo worktree created", prefix, false, true)
	}
	// Worktree stage all
	if p.Err == nil {
		addOpt := git.AddOptions{
			All:  true,
			Path: gitDir,
		}
		p.Err = gitWorktree.AddWithOptions(&addOpt)
		helper.ReportDebug("worktree staged("+gitDir+")", prefix, false, true)
	}
	// Worktree commit
	if p.Err == nil {
		commit, p.Err = gitWorktree.Commit(msg, &commitOptions)
		// commit, p.Err = gitWorktree.Commit(msg, nil)
		helper.ReportDebug("worktree committed", prefix, false, true)
	}
	// Repository commit
	if p.Err == nil {
		// commitObj, p.Err = p.gitRepo.CommitObject(commit)
		_, p.Err = gitRepo.CommitObject(commit)
		helper.ReportDebug("repo committed", prefix, false, true)
	}
	// Repository tag
	if tag {
		if p.Err == nil {
			gitHead, p.Err = gitRepo.Head()
		}
		if p.Err == nil {
			_, p.Err = gitRepo.CreateTag(msg, gitHead.Hash(), nil)
			helper.ReportDebug("tag("+msg+")", prefix, false, true)
		}
	}
	helper.ErrsQueue(p.Err, prefix)
	return p
}

func (p *TypeRepository) copyDir(dirSrc string, dirDest string) *TypeRepository {
	prefix := p.myType + ".copyDir"
	helper.ReportDebug("-- Start", prefix, false, true)
	if p.Err != nil {
		return p
	}
	if p.Err == nil {
		os.RemoveAll(dirDest) // Delete first
		if p.Err = os.CopyFS(dirDest, os.DirFS(dirSrc)); p.Err != nil {
			helper.ErrsQueue(p.Err, prefix)
		}
	}
	return p
}
