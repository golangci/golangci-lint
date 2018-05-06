package golinters

import (
	"context"
	"fmt"
	"go/build"
	"io/ioutil"
	"path/filepath"

	"github.com/golang/lint"
	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/result"
	"github.com/golangci/golangci-shared/pkg/executors"
)

type golint struct{}

func (golint) Name() string {
	return "golint"
}

func (g golint) Run(ctx context.Context, exec executors.Executor, cfg *config.Run) (*result.Result, error) {
	var issues []result.Issue
	if cfg.Paths.IsDirsRun {
		for _, path := range cfg.Paths.Dirs {
			i, err := lintDir(path, cfg.Golint.MinConfidence)
			if err != nil {
				// TODO: skip and warn
				return nil, fmt.Errorf("can't lint dir %s: %s", path, err)
			}
			issues = append(issues, i...)
		}
	} else {
		i, err := lintFiles(cfg.Golint.MinConfidence, cfg.Paths.Files...)
		if err != nil {
			// TODO: skip and warn
			return nil, fmt.Errorf("can't lint files %s: %s", cfg.Paths.Files, err)
		}
		issues = append(issues, i...)
	}

	return &result.Result{
		Issues: issues,
	}, nil
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
				File:       p.Position.Filename,
				LineNumber: p.Position.Line,
				Text:       p.Text,
			})
			// TODO: use p.Link and p.Category
		}
	}

	return issues, nil
}
