package golinters

import (
	"github.com/bkielbasa/cyclop/pkg/analyzer"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/goanalysis"
)

func NewCyclop(settings *config.Cyclop) *goanalysis.Linter {
	a := analyzer.NewAnalyzer()

	var cfg map[string]map[string]any
	if settings != nil {
		d := map[string]any{
			"skipTests": settings.SkipTests,
		}

		if settings.MaxComplexity != 0 {
			d["maxComplexity"] = settings.MaxComplexity
		}

		if settings.PackageAverage != 0 {
			d["packageAverage"] = settings.PackageAverage
		}

		cfg = map[string]map[string]any{a.Name: d}
	}

	return goanalysis.NewLinter(
		a.Name,
		"checks function and package cyclomatic complexity",
		[]*analysis.Analyzer{a},
		cfg,
	).WithLoadMode(goanalysis.LoadModeTypesInfo)
}
