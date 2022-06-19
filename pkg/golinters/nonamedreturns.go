package golinters

import (
	"github.com/firefart/nonamedreturns/analyzer"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"
)

func NewNoNamedReturns(settings *config.NoNamedReturnsSettings) *goanalysis.Linter {
	a := analyzer.Analyzer

	var cfg map[string]map[string]interface{}
	if settings != nil {
		cfg = map[string]map[string]interface{}{
			a.Name: {
				analyzer.FlagAllowErrorInDefer: settings.AllowErrorInDefer,
			},
		}
	}

	return goanalysis.NewLinter(
		a.Name,
		a.Doc,
		[]*analysis.Analyzer{a},
		cfg,
	).WithLoadMode(goanalysis.LoadModeSyntax)
}
