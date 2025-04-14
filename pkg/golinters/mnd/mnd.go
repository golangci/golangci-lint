package mnd

import (
	mnd "github.com/tommy-muehle/go-mnd/v2"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/v2/pkg/config"
	"github.com/golangci/golangci-lint/v2/pkg/goanalysis"
)

func New(settings *config.MndSettings) *goanalysis.Linter {
	a := mnd.Analyzer

	var linterCfg map[string]map[string]any
	if settings != nil {
		cfg := make(map[string]any)
		if len(settings.Checks) > 0 {
			cfg["checks"] = settings.Checks
		}
		if len(settings.IgnoredNumbers) > 0 {
			cfg["ignored-numbers"] = settings.IgnoredNumbers
		}
		if len(settings.IgnoredFiles) > 0 {
			cfg["ignored-files"] = settings.IgnoredFiles
		}
		if len(settings.IgnoredFunctions) > 0 {
			cfg["ignored-functions"] = settings.IgnoredFunctions
		}

		linterCfg = map[string]map[string]any{
			a.Name: cfg,
		}
	}

	return goanalysis.NewLinter(
		a.Name,
		"An analyzer to detect magic numbers.",
		[]*analysis.Analyzer{a},
		linterCfg,
	).WithLoadMode(goanalysis.LoadModeSyntax)
}
