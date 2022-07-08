package golinters

import (
	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"

	"github.com/sivchari/tenv"
	"golang.org/x/tools/go/analysis"
)

func NewTenv(settings *config.TenvSettings) *goanalysis.Linter {
	a := tenv.Analyzer

	var cfg map[string]map[string]interface{}
	if settings != nil {
		cfg = map[string]map[string]interface{}{
			a.Name: {
				tenv.A: settings.All,
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
