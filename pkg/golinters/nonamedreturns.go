package golinters

import (
	"github.com/firefart/nonamedreturns/analyzer"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/goanalysis"
)

func NewNoNamedReturns(settings *config.NoNamedReturnsSettings) *goanalysis.Linter {
	a := analyzer.Analyzer

	var cfg map[string]map[string]any
	if settings != nil {
		cfg = map[string]map[string]any{
			a.Name: {
				analyzer.FlagReportErrorInDefer: settings.ReportErrorInDefer,
			},
		}
	}

	return goanalysis.NewLinter(
		a.Name,
		a.Doc,
		[]*analysis.Analyzer{a},
		cfg,
	).WithLoadMode(goanalysis.LoadModeTypesInfo)
}
