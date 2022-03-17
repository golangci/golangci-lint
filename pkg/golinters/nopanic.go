package golinters

import (
	"github.com/vladopajic/nopanic/pkg/nopanic"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"
)

func NewNoPanic(settings *config.NoPanicSettings) *goanalysis.Linter {
	analyzers := []*analysis.Analyzer{nopanic.NewAnalyzer()}

	cfg := map[string]map[string]interface{}{}
	if settings != nil {
		cfg[nopanic.AnalyzerName] = map[string]interface{}{
			nopanic.FlagAllowPanicMainFunc:    settings.AllowPanicMainFunc,
			nopanic.FlagAllowPanicMainPackage: settings.AllowPanicMainPackage,
			nopanic.FlagAllowExitMainFunc:     settings.AllowExitMainFunc,
			nopanic.FlagAllowExitMainPackage:  settings.AllowExitMainPackage,
		}
	}

	return goanalysis.NewLinter(
		nopanic.AnalyzerName,
		nopanic.AnalyzerDoc,
		analyzers,
		cfg,
	).WithLoadMode(goanalysis.LoadModeTypesInfo)
}
