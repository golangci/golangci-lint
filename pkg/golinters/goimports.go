package golinters

import (
	"sync"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/imports"

	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"
	"github.com/golangci/golangci-lint/pkg/lint/linter"
)

const goimportsName = "goimports"

func NewGoimports() *goanalysis.Linter {
	var (
		mu        sync.Mutex
		resIssues []goanalysis.Issue
		options   = &imports.Options{
			TabWidth:  8,
			TabIndent: true,
			Comments:  true,
			Fragment:  true,
		}
	)

	analyzer := &analysis.Analyzer{
		Name: goimportsName,
		Doc:  goanalysis.TheOnlyanalyzerDoc,
	}
	return goanalysis.NewLinter(
		goimportsName,
		"In addition to fixing imports, goimports also formats your code in the same style as gofmt.",
		[]*analysis.Analyzer{analyzer},
		nil,
	).WithContextSetter(func(lintCtx *linter.Context) {
		imports.LocalPrefix = lintCtx.Settings().Goimports.LocalPrefixes
		analyzer.Run = func(pass *analysis.Pass) (interface{}, error) {
			cb := func(filename string, src []byte) ([]byte, error) {
				return imports.Process(filename, src, options)
			}

			issues, err := runFormatAndDiffLinter(pass, lintCtx, goimportsName, cb)
			if err != nil {
				return nil, err
			}
			if len(issues) == 0 {
				return nil, nil
			}

			mu.Lock()
			resIssues = append(resIssues, issues...)
			mu.Unlock()

			return nil, nil
		}
	}).WithIssuesReporter(func(*linter.Context) []goanalysis.Issue {
		return resIssues
	}).WithLoadMode(goanalysis.LoadModeSyntax)
}
