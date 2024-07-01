package fparams

import (
	"github.com/artemk1337/fparams/pkg/analyzer"
	"github.com/golangci/golangci-lint/pkg/config"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/goanalysis"
)

func New(settings *config.Fparams) *goanalysis.Linter {
	a := analyzer.NewAnalyzer()

	cfg := map[string]map[string]any{}
	if settings != nil {
		cfg[a.Name] = map[string]any{
			"disableCheckFuncParams":  settings.DisableCheckFuncParams,
			"disableCheckFuncReturns": settings.DisableCheckFuncReturns,
		}
	}

	return goanalysis.NewLinter(
		a.Name,
		a.Doc,
		[]*analysis.Analyzer{a},
		nil,
	).WithLoadMode(goanalysis.LoadModeTypesInfo)
}
