package contextcheck

import (
	"github.com/kkHAIKE/contextcheck"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/v2/pkg/goanalysis"
	"github.com/golangci/golangci-lint/v2/pkg/lint/linter"
)

func New() *goanalysis.Linter {
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
