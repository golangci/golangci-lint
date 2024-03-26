package golinters

import (
	"github.com/macabu/inamedparam"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/goanalysis"
)

func NewINamedParam(settings *config.INamedParamSettings) *goanalysis.Linter {
	a := inamedparam.Analyzer

	var cfg map[string]map[string]any

	if settings != nil {
		cfg = map[string]map[string]any{
			a.Name: {
				"skip-single-param": settings.SkipSingleParam,
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
