package golinters

import (
	"context"
	"fmt"
	"io/ioutil"

	"github.com/golang/lint"
	"github.com/golangci/golangci-lint/pkg/result"
)

type Golint struct{}

func (Golint) Name() string {
	return "golint"
}

func (g Golint) Run(ctx context.Context, lintCtx *Context) ([]result.Issue, error) {
	var issues []result.Issue
	for _, pkgFiles := range lintCtx.Paths.FilesGrouppedByDirs() {
		i, err := lintFiles(lintCtx.RunCfg().Golint.MinConfidence, pkgFiles...)
		if err != nil {
			// TODO: skip and warn
			return nil, fmt.Errorf("can't lint files %s: %s", lintCtx.Paths.Files, err)
		}
		issues = append(issues, i...)
	}

	return issues, nil
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
