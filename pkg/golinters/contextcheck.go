package golinters

import (
	"github.com/sylvia7788/contextcheck"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"
	"github.com/golangci/golangci-lint/pkg/lint/linter"
)

func NewContextCheck() *goanalysis.Linter {
	analyzer := contextcheck.NewAnalyzer()
	return goanalysis.NewLinter(
		"contextcheck",
		"check the function whether use a non-inherited context",
		[]*analysis.Analyzer{analyzer},
		nil,
	).WithLoadMode(goanalysis.LoadModeTypesInfo).WithNoCache().
		WithContextSetter(func(lintCtx *linter.Context) {
			analyzer.Run = contextcheck.NewRun(lintCtx.Packages)
		})
}
