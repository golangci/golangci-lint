package golinters

import (
	"strings"

	"github.com/mbilski/exhaustivestruct/pkg/analyzer"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"
)

func NewExhaustiveStruct(settings *config.ExhaustiveStructSettings) *goanalysis.Linter {
	a := analyzer.Analyzer

	var cfg map[string]map[string]any
	if settings != nil {
		cfg = map[string]map[string]any{
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
