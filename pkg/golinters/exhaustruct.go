package golinters

import (
	"strings"

	"github.com/GaijinEntertainment/go-exhaustruct/pkg/analyzer"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"
)

func NewExhaustruct(settings *config.ExhaustructSettings) *goanalysis.Linter {
	a := analyzer.Analyzer

	var cfg map[string]map[string]interface{}
	if settings != nil {
		cfg = map[string]map[string]interface{}{
			a.Name: {
				"include": strings.Join(settings.Include, ","),
				"exclude": strings.Join(settings.Exclude, ","),
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
