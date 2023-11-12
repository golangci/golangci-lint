package golinters

import (
	"github.com/maranqz/go-factory-lint"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"
)

func NewGoFactoryLint(settings *config.GoFactoryLintSettings) *goanalysis.Linter {
	a := factory.NewAnalyzer()

	cfg := make(map[string]map[string]any)
	if settings != nil {
		cfg[a.Name] = map[string]any{}

		if len(settings.BlockedPkgs) > 0 {
			cfg[a.Name]["blockedPkgs"] = settings.BlockedPkgs
			cfg[a.Name]["onlyBlockedPkgs"] = settings.OnlyBlockedPkgs
		}
	}

	return goanalysis.NewLinter(
		a.Name,
		a.Doc,
		[]*analysis.Analyzer{a},
		cfg,
	).WithLoadMode(goanalysis.LoadModeTypesInfo)
}
