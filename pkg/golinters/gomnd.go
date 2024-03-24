package golinters

import (
	mnd "github.com/tommy-muehle/go-mnd/v2"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/goanalysis"
)

func NewGoMND(settings *config.GoMndSettings) *goanalysis.Linter {
	// The constant is only used to force the analyzer name to use the same name as the linter.
	// This is required to avoid displaying the analyzer name inside the issue text.
	//
	// Alternative names cannot help here because of the linter configuration that uses `gomnd` as a name.
	// The complexity of handling alternative names at a lower level (i.e. `goanalysis.Linter`) isn't worth the cost.
	// The only way to handle it properly is to deprecate and "duplicate" the linter and its configuration,
	// for now, I don't know if it's worth the cost.
	// TODO(ldez): in v2, rename to mnd as the real analyzer name?
	const name = "gomnd"

	a := mnd.Analyzer
	a.Name = name

	var linterCfg map[string]map[string]any

	if settings != nil {
		// Convert deprecated setting.
		if len(settings.Settings) > 0 {
			linterCfg = settings.Settings
		} else {
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
	}

	return goanalysis.NewLinter(
		a.Name,
		"An analyzer to detect magic numbers.",
		[]*analysis.Analyzer{a},
		linterCfg,
	).WithLoadMode(goanalysis.LoadModeSyntax)
}
