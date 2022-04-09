package golinters

import (
	"github.com/mweb/floatcompare"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"
)

func NewFloatCompare(settings *config.FloatCompareSettings) *goanalysis.Linter {
	a := floatcompare.NewAnalyzer()

	var cfg map[string]map[string]interface{}
	if settings != nil {
		d := map[string]interface{}{
			"equalOnly": settings.EqualOnly,
		}

		cfg = map[string]map[string]interface{}{a.Name: d}
	}

	return goanalysis.NewLinter(
		a.Name,
		a.Doc,
		[]*analysis.Analyzer{a},
		cfg,
	).WithLoadMode(goanalysis.LoadModeTypesInfo)
}
