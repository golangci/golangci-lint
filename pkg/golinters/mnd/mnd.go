package mnd

import (
	mnd "github.com/tommy-muehle/go-mnd/v2"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/goanalysis"
)

func New(settings *config.MndSettings) *goanalysis.Linter {
	return newMND(mnd.Analyzer, settings, nil)
}

func NewGoMND(settings *config.GoMndSettings) *goanalysis.Linter {
	// shallow copy because mnd.Analyzer is a global variable.
	a := new(analysis.Analyzer)
	*a = *mnd.Analyzer

	// Used to force the analyzer name to use the same name as the linter.
	// This is required to avoid displaying the analyzer name inside the issue text.
	a.Name = "gomnd"

	var linterCfg map[string]map[string]any

	if settings != nil && len(settings.Settings) > 0 {
		// Convert deprecated setting.
		linterCfg = map[string]map[string]any{
			a.Name: settings.Settings["mnd"],
		}
	}

	return newMND(a, &settings.MndSettings, linterCfg)
}

func newMND(a *analysis.Analyzer, settings *config.MndSettings, linterCfg map[string]map[string]any) *goanalysis.Linter {
	if len(linterCfg) == 0 && settings != nil {
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
