package contextcheck

import (
	"github.com/kkHAIKE/contextcheck"

	"github.com/golangci/golangci-lint/v2/pkg/goanalysis"
	"github.com/golangci/golangci-lint/v2/pkg/lint/linter"
)

func New() *goanalysis.Linter {
	analyzer := contextcheck.NewAnalyzer(contextcheck.Configuration{})

	return goanalysis.
		NewLinterFromAnalyzer(analyzer).
		WithContextSetter(func(lintCtx *linter.Context) {
			analyzer.Run = contextcheck.NewRun(lintCtx.Packages, false)
		}).
		WithLoadMode(goanalysis.LoadModeTypesInfo)
}
