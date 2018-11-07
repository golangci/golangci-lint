package golinters

import (
	"context"
	"fmt"
	"go/ast"
	"go/types"
	"path/filepath"
	"runtime"
	"runtime/debug"
	"strings"
	"sync"

	_ "github.com/go-critic/checkers" // this import register checkers

	"github.com/go-lintpack/lintpack"
	"golang.org/x/tools/go/loader"

	"github.com/golangci/golangci-lint/pkg/lint/linter"
	"github.com/golangci/golangci-lint/pkg/result"
)

type Gocritic struct{}

func (Gocritic) Name() string {
	return "gocritic"
}

func (Gocritic) Desc() string {
	return "The most opinionated Go source code linter"
}

func (lint Gocritic) Run(ctx context.Context, lintCtx *linter.Context) ([]result.Issue, error) {
	sizes := types.SizesFor("gc", runtime.GOARCH)
	lintpackCtx := lintpack.NewContext(lintCtx.Program.Fset, sizes)

	s := lintCtx.Settings().Gocritic
	var enabledCheckers []*lintpack.Checker
	for _, info := range lintpack.GetCheckersInfo() {
		if !s.IsCheckEnabled(info.Name) {
			continue
		}

		params := s.SettingsPerCheck[strings.ToLower(info.Name)]
		c := lintpack.NewChecker(lintpackCtx, info, params)
		enabledCheckers = append(enabledCheckers, c)
	}

	issuesCh := make(chan result.Issue, 1024)
	var panicErr error
	go func() {
		defer func() {
			if err := recover(); err != nil {
				panicErr = fmt.Errorf("panic occured: %s", err)
				lintCtx.Log.Warnf("Panic: %s", debug.Stack())
			}
		}()

		for _, pkgInfo := range lintCtx.Program.InitialPackages() {
			lintpackCtx.SetPackageInfo(&pkgInfo.Info, pkgInfo.Pkg)
			lint.runOnPackage(lintpackCtx, enabledCheckers, pkgInfo, issuesCh)
		}
		close(issuesCh)
	}()

	var res []result.Issue
	for i := range issuesCh {
		res = append(res, i)
	}
	if panicErr != nil {
		return nil, panicErr
	}

	return res, nil
}

func (lint Gocritic) runOnPackage(lintpackCtx *lintpack.Context, checkers []*lintpack.Checker,
	pkgInfo *loader.PackageInfo, ret chan<- result.Issue) {

	for _, f := range pkgInfo.Files {
		filename := filepath.Base(lintpackCtx.FileSet.Position(f.Pos()).Filename)
		lintpackCtx.SetFileInfo(filename, f)

		lint.runOnFile(lintpackCtx, f, checkers, ret)
	}
}

func (lint Gocritic) runOnFile(ctx *lintpack.Context, f *ast.File, checkers []*lintpack.Checker,
	ret chan<- result.Issue) {

	var wg sync.WaitGroup
	wg.Add(len(checkers))
	for _, c := range checkers {
		// All checkers are expected to use *lint.Context
		// as read-only structure, so no copying is required.
		go func(c *lintpack.Checker) {
			defer wg.Done()

			for _, warn := range c.Check(f) {
				pos := ctx.FileSet.Position(warn.Node.Pos())
				ret <- result.Issue{
					Pos:        pos,
					Text:       fmt.Sprintf("%s: %s", c.Info.Name, warn.Text),
					FromLinter: lint.Name(),
				}
			}
		}(c)
	}

	wg.Wait()
}
