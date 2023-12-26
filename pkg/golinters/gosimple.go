package golinters

import (
	"honnef.co/go/tools/simple"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"
)

func NewGosimple(settings *config.CommonStaticCheckSettings) *goanalysis.Linter {
	cfg := commonStaticCheckConfig(settings)

	analyzers := setupStaticCheckAnalyzers(simple.Analyzers, settings.GetGoVersion(), cfg.Checks)

	return goanalysis.NewLinter(
		"gosimple",
		"Linter for Go source code that specializes in simplifying code",
		analyzers,
		nil,
	).WithLoadMode(goanalysis.LoadModeTypesInfo)
}
