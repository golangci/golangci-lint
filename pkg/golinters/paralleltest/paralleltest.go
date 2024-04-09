package paralleltest

import (
	"github.com/kunwardeep/paralleltest/pkg/paralleltest"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/goanalysis"
)

func New(settings *config.ParallelTestSettings) *goanalysis.Linter {
	a := paralleltest.NewAnalyzer()

	var cfg map[string]map[string]any
	if settings != nil {
		d := map[string]any{
			"i":                     settings.IgnoreMissing,
			"ignoremissingsubtests": settings.IgnoreMissingSubtests,
		}

		if config.IsGoGreaterThanOrEqual(settings.Go, "1.22") {
			d["ignoreloopVar"] = true
		}

		cfg = map[string]map[string]any{a.Name: d}
	}

	return goanalysis.NewLinter(
		a.Name,
		"Detects missing usage of t.Parallel() method in your Go test",
		[]*analysis.Analyzer{a},
		cfg,
	).WithLoadMode(goanalysis.LoadModeTypesInfo)
}
