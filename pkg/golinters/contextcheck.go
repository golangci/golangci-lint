package golinters

import (
	"github.com/sylvia7788/contextcheck"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"
	"github.com/golangci/golangci-lint/pkg/lint/linter"
)

func NewContextCheck(settings *config.ContextCheckSettings) *goanalysis.Linter {
	conf := contextcheck.Configuration{}
	if settings != nil {
		conf.DisableFact = settings.DisableFact
	}
	analyzer := contextcheck.NewAnalyzer(conf)
	return goanalysis.NewLinter(
		"contextcheck",
		"check the function whether use a non-inherited context",
		[]*analysis.Analyzer{analyzer},
		nil,
	).WithLoadMode(goanalysis.LoadModeTypesInfo).
		WithContextSetter(func(lintCtx *linter.Context) {
			analyzer.Run = contextcheck.NewRun(lintCtx.Packages, conf.DisableFact)
		})
}
