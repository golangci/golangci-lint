package golinters

import (
	"github.com/alingse/asasalint"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"
)

func NewAsasalint(setting *config.AsasalintSettings) *goanalysis.Linter {
	cfg := asasalint.LinterSetting{}
	if setting != nil {
		cfg.Exclude = setting.Exclude
		cfg.NoBuiltinExclusions = !setting.UseBuiltinExclusions
		cfg.IgnoreTest = setting.IgnoreTest
	}

	a, err := asasalint.NewAnalyzer(cfg)
	if err != nil {
		linterLogger.Fatalf("asasalint: create analyzer: %v", err)
	}

	return goanalysis.NewLinter(
		a.Name,
		a.Doc,
		[]*analysis.Analyzer{a},
		nil,
	).WithLoadMode(goanalysis.LoadModeTypesInfo)
}
