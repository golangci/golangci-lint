package perfsprint

import (
	"github.com/catenacyber/perfsprint/analyzer"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/goanalysis"
)

func New(settings *config.PerfSprintSettings) *goanalysis.Linter {
	a := analyzer.New()

	cfg := map[string]map[string]any{
		a.Name: {"fiximports": false},
	}

	if settings != nil {
		cfg[a.Name]["int-conversion"] = settings.IntConversion
		cfg[a.Name]["err-error"] = settings.ErrError
		cfg[a.Name]["errorf"] = settings.ErrorF
		cfg[a.Name]["sprintf1"] = settings.SprintF1
		cfg[a.Name]["strconcat"] = settings.StrConcat
	}

	return goanalysis.NewLinter(
		a.Name,
		a.Doc,
		[]*analysis.Analyzer{a},
		cfg,
	).WithLoadMode(goanalysis.LoadModeTypesInfo)
}
