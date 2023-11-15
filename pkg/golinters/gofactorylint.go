package golinters

import (
	"github.com/maranqz/go-factory-lint"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"
)

func NewGoFactoryLint(settings *config.GoFactoryLintSettings) *goanalysis.Linter {
	analyzer := factory.NewAnalyzer()

	cfg := make(map[string]map[string]any)
	if settings != nil {
		cfg[analyzer.Name] = map[string]any{}

		if len(settings.BlockedPkgs) > 0 {
			cfg[analyzer.Name]["blockedPkgs"] = settings.BlockedPkgs
			cfg[analyzer.Name]["onlyBlockedPkgs"] = settings.OnlyBlockedPkgs
		}
	}

	return goanalysis.NewLinter(
		analyzer.Name,
		analyzer.Doc,
		[]*analysis.Analyzer{analyzer},
		cfg,
	).WithLoadMode(goanalysis.LoadModeTypesInfo)
}
