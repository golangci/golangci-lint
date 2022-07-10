package golinters

import (
	"github.com/alingse/asasalint"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"
)

func NewAsasalint(cfg *config.AsasalintSettings) *goanalysis.Linter {
	setting := asasalint.LinterSetting{}
	if cfg != nil {
		setting.Exclude = cfg.Exclude
		setting.NoDefaultExclude = cfg.NoDefaultExclude
		setting.IgnoreInTest = cfg.IgnoreInTest
	}
	a := asasalint.NewAnalyzer(setting)

	return goanalysis.NewLinter(
		a.Name,
		a.Doc,
		[]*analysis.Analyzer{a},
		nil,
	).WithLoadMode(goanalysis.LoadModeTypesInfo)
}
