package ttempdir

import (
	"github.com/peczenyj/ttempdir"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/goanalysis"
)

func New(settings *config.TtempdirSettings) *goanalysis.Linter {
	a := ttempdir.Analyzer

	var cfg map[string]map[string]any
	if settings != nil {
		cfg = map[string]map[string]any{
			a.Name: {
				ttempdir.A:   settings.All,
				ttempdir.MRL: settings.MaxRecursionLevel,
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
