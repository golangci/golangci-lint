package golinters

import (
	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"
	"github.com/raffepaffe/nevernester/pkg/analyzer"
	"golang.org/x/tools/go/analysis"
)

const nevernesterName = "nevernester"

func NewNeverNester(settings *config.NeverNesterSettings) *goanalysis.Linter {
	a := analyzer.New()

	var cfg map[string]map[string]interface{}
	if settings != nil {
		d := map[string]interface{}{
			"skip-tests":      settings.SkipTests,
			"skip-benchmarks": settings.SkipBenchmarks,
		}

		if settings.MaxNesting != 0 {
			d["max-nesting"] = settings.MaxNesting
		}

		cfg = map[string]map[string]interface{}{a.Name: d}
	}

	return goanalysis.NewLinter(
		nevernesterName,
		"checks nesting in functions",
		[]*analysis.Analyzer{a},
		cfg,
	).WithLoadMode(goanalysis.LoadModeTypesInfo)
}
