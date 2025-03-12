package cyclop

import (
	"github.com/bkielbasa/cyclop/pkg/analyzer"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/v2/pkg/config"
	"github.com/golangci/golangci-lint/v2/pkg/goanalysis"
)

func New(settings *config.CyclopSettings) *goanalysis.Linter {
	a := analyzer.NewAnalyzer()

	var cfg map[string]map[string]any
	if settings != nil {
		d := map[string]any{
			// Should be managed with `linters.exclusions.rules`.
			"skipTests": false,
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
		a.Doc,
		[]*analysis.Analyzer{a},
		cfg,
	).WithLoadMode(goanalysis.LoadModeSyntax)
}
