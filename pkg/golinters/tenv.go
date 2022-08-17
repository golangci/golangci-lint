package golinters

import (
	"github.com/sivchari/tenv"
	"golang.org/x/tools/go/analysis"

	"github.com/anduril/golangci-lint/pkg/config"
	"github.com/anduril/golangci-lint/pkg/golinters/goanalysis"
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
