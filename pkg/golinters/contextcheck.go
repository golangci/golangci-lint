package golinters

import (
	"github.com/sylvia7788/contextcheck"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"
	"github.com/golangci/golangci-lint/pkg/lint/linter"
)

func NewContextCheck() *goanalysis.Linter {
	analyzer := contextcheck.NewAnalyzer(contextcheck.Configuration{})

	return goanalysis.NewLinter(
		analyzer.Name,
		analyzer.Doc,
		[]*analysis.Analyzer{analyzer},
		nil,
	).WithContextSetter(func(lintCtx *linter.Context) {
		analyzer.Run = contextcheck.NewRun(lintCtx.Packages, false)
	}).WithLoadMode(goanalysis.LoadModeTypesInfo)
}
