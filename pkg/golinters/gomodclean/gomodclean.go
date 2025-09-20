package gomodclean

import (
	"sync"

	gomodcleananalyzer "github.com/dmrioja/gomodclean/pkg/analyzer"

	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/v2/pkg/goanalysis"
	"github.com/golangci/golangci-lint/v2/pkg/lint/linter"
	"github.com/golangci/golangci-lint/v2/pkg/result"
)

const linterName = "gomodclean"

func New() *goanalysis.Linter {
	var issues []*goanalysis.Issue
	var once sync.Once

	analyzer := &analysis.Analyzer{
		Name: linterName,
		Doc:  "Linter to check dependencies are well structured inside your go.mod file.",
		Run:  goanalysis.DummyRun,
	}

	return goanalysis.
		NewLinterFromAnalyzer(analyzer).
		WithContextSetter(func(lintCtx *linter.Context) {
			analyzer.Run = func(pass *analysis.Pass) (any, error) {
				once.Do(func() {
					results, err := gomodcleananalyzer.Analyze()
					if err != nil {
						lintCtx.Log.Warnf("running %s failed: %s: "+
							"if you are not using go modules it is suggested to disable this linter", linterName, err)
						return
					}

					for _, p := range results {
						issues = append(issues, goanalysis.NewIssue(&result.Issue{
							FromLinter: linterName,
							Pos:        p.Position,
							Text:       p.Text,
						}, pass))
					}
				})

				return nil, nil
			}
		}).
		WithIssuesReporter(func(*linter.Context) []*goanalysis.Issue {
			return issues
		}).
		WithLoadMode(goanalysis.LoadModeSyntax)
}
