package golinters

import (
	"context"
	"fmt"
	"go/build"
	"io/ioutil"
	"path/filepath"

	"github.com/golang/lint"
	"github.com/golangci/golangci-lint/pkg/result"
)

type Golint struct{}

func (Golint) Name() string {
	return "golint"
}

func (g Golint) Run(ctx context.Context, lintCtx *Context) ([]result.Issue, error) {
	var issues []result.Issue
	if lintCtx.Paths.IsDirsRun {
		for _, path := range lintCtx.Paths.Dirs {
			i, err := lintDir(path, lintCtx.RunCfg().Golint.MinConfidence)
			if err != nil {
				// TODO: skip and warn
				return nil, fmt.Errorf("can't lint dir %s: %s", path, err)
			}
			issues = append(issues, i...)
		}
	} else {
		i, err := lintFiles(lintCtx.RunCfg().Golint.MinConfidence, lintCtx.Paths.Files...)
		if err != nil {
			// TODO: skip and warn
			return nil, fmt.Errorf("can't lint files %s: %s", lintCtx.Paths.Files, err)
		}
		issues = append(issues, i...)
	}

	return issues, nil
}

func lintDir(dirname string, minConfidence float64) ([]result.Issue, error) {
	pkg, err := build.ImportDir(dirname, 0)
	if err != nil {
		if _, nogo := err.(*build.NoGoError); nogo {
			// Don't complain if the failure is due to no Go source files.
			return nil, nil
		}

		return nil, fmt.Errorf("can't import dir %s", dirname)
	}

	return lintImportedPackage(pkg, minConfidence)
}

func lintImportedPackage(pkg *build.Package, minConfidence float64) ([]result.Issue, error) {
	var files []string
	files = append(files, pkg.GoFiles...)
	files = append(files, pkg.CgoFiles...)
	files = append(files, pkg.TestGoFiles...)
	files = append(files, pkg.XTestGoFiles...)
	if pkg.Dir != "." {
		for i, f := range files {
			files[i] = filepath.Join(pkg.Dir, f)
		}
	}

	return lintFiles(minConfidence, files...)
}

func lintFiles(minConfidence float64, filenames ...string) ([]result.Issue, error) {
	files := make(map[string][]byte)
	for _, filename := range filenames {
		src, err := ioutil.ReadFile(filename)
		if err != nil {
			return nil, fmt.Errorf("can't read file %s: %s", filename, err)
		}
		files[filename] = src
	}

	l := new(lint.Linter)
	ps, err := l.LintFiles(files)
	if err != nil {
		return nil, fmt.Errorf("can't lint files %s: %s", filenames, err)
	}

	var issues []result.Issue
	for _, p := range ps {
		if p.Confidence >= minConfidence {
			issues = append(issues, result.Issue{
				Pos:  p.Position,
				Text: p.Text,
			})
			// TODO: use p.Link and p.Category
		}
	}

	return issues, nil
}
