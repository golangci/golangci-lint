package golinters

import (
	mnd "github.com/tommy-muehle/go-mnd/v2"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"
)

func NewGoMND(settings *config.GoMndSettings) *goanalysis.Linter {
	var linterCfg map[string]map[string]interface{}

	if settings != nil {
		// TODO(ldez) For compatibility only, must be drop in v2.
		if len(settings.Settings) > 0 {
			linterCfg = settings.Settings
		} else {
			linterCfg = map[string]map[string]interface{}{
				"mnd": {
					"checks":            settings.Checks,
					"ignored-numbers":   settings.IgnoredNumbers,
					"ignored-files":     settings.IgnoredFiles,
					"ignored-functions": settings.IgnoredFunctions,
				},
			}
		}
	}

	return goanalysis.NewLinter(
		"gomnd",
		"An analyzer to detect magic numbers.",
		[]*analysis.Analyzer{mnd.Analyzer},
		linterCfg,
	).WithLoadMode(goanalysis.LoadModeSyntax)
}
