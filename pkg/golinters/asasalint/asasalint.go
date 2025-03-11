package asasalint

import (
	"github.com/alingse/asasalint"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/v2/pkg/config"
	"github.com/golangci/golangci-lint/v2/pkg/goanalysis"
	"github.com/golangci/golangci-lint/v2/pkg/golinters/internal"
)

func New(settings *config.AsasalintSettings) *goanalysis.Linter {
	cfg := asasalint.LinterSetting{}
	if settings != nil {
		cfg.Exclude = settings.Exclude
		cfg.NoBuiltinExclusions = !settings.UseBuiltinExclusions

		// Should be managed with `linters.exclusions.rules`.
		cfg.IgnoreTest = false
	}

	a, err := asasalint.NewAnalyzer(cfg)
	if err != nil {
		internal.LinterLogger.Fatalf("asasalint: create analyzer: %v", err)
	}

	return goanalysis.NewLinter(
		a.Name,
		a.Doc,
		[]*analysis.Analyzer{a},
		nil,
	).WithLoadMode(goanalysis.LoadModeTypesInfo)
}
