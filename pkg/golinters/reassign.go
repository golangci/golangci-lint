package golinters

import (
	"github.com/curioswitch/go-reassign"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"
)

func NewReassign(settings *config.ReassignSettings) *goanalysis.Linter {
	a := reassign.NewAnalyzer()

	var cfg map[string]map[string]interface{}
	if settings != nil {
		cfg = map[string]map[string]interface{}{
			a.Name: {
				reassign.FlagPattern: settings.Pattern,
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
