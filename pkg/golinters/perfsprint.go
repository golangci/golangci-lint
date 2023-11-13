package golinters

import (
	"github.com/catenacyber/perfsprint/analyzer"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"
)

func NewPerfSprint(settings *config.PerfSprintSettings) *goanalysis.Linter {
	a := analyzer.New()

	var cfg map[string]map[string]any
	if settings != nil {
		cfg = map[string]map[string]any{
			a.Name: {
				"int-conversion": settings.IntConversion,
				"err-error":      settings.ErrError,
				"errorf":         settings.ErrorF,
				"sprintf1":       settings.SprintF1,
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
