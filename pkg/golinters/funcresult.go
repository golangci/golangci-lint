package golinters

import (
	funcresult "github.com/leonklingele/funcresult/pkg/analyzer"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"
)

func NewFuncResult(settings *config.FuncResultSettings) *goanalysis.Linter {
	linterCfg := map[string]map[string]interface{}{}
	if settings != nil {
		linterCfg["funcresult"] = map[string]interface{}{
			"require-named":   settings.RequireNamed,
			"require-unnamed": settings.RequireUnnamed,
		}
	}

	return goanalysis.NewLinter(
		"funcresult",
		"An analyzer to function result parameters.",
		[]*analysis.Analyzer{funcresult.New()},
		linterCfg,
	).WithLoadMode(goanalysis.LoadModeSyntax)
}
