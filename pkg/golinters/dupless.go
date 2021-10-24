package golinters

import (
	"github.com/bkielbasa/dupless/pkg/analyzer"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"
)

const duplessName = "dupless"

func NewDupless(settings *config.DuplessSettings) *goanalysis.Linter {
	a := analyzer.NewAnalyzer()

	var cfg map[string]map[string]interface{}
	if settings != nil {
		d := map[string]interface{}{}

		if len(settings.PackageNames) != 0 {
			d["packageNames"] = settings.PackageNames
		}

		if len(settings.FunctionNames) != 0 {
			d["functionNames"] = settings.FunctionNames
		}

		if len(settings.VariableNames) != 0 {
			d["variableNames"] = settings.VariableNames
		}

		cfg = map[string]map[string]interface{}{a.Name: d}
	}
	return goanalysis.NewLinter(
		duplessName,
		"check if functions, packages or variables contain forbidden patterns",
		[]*analysis.Analyzer{
			analyzer.NewAnalyzer(),
		},
		cfg,
	).WithLoadMode(goanalysis.LoadModeTypesInfo)
}
