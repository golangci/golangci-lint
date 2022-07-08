package golinters

import (
	"strings"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"

	"github.com/mbilski/exhaustivestruct/pkg/analyzer"
	"golang.org/x/tools/go/analysis"
)

func NewExhaustiveStruct(settings *config.ExhaustiveStructSettings) *goanalysis.Linter {
	a := analyzer.Analyzer

	var cfg map[string]map[string]interface{}
	if settings != nil {
		cfg = map[string]map[string]interface{}{
			a.Name: {
				"struct_patterns": strings.Join(settings.StructPatterns, ","),
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
