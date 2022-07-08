package golinters

import (
	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"

	"honnef.co/go/tools/simple"
)

func NewGosimple(settings *config.StaticCheckSettings) *goanalysis.Linter {
	cfg := staticCheckConfig(settings)

	analyzers := setupStaticCheckAnalyzers(simple.Analyzers, getGoVersion(settings), cfg.Checks)

	return goanalysis.NewLinter(
		"gosimple",
		"Linter for Go source code that specializes in simplifying code",
		analyzers,
		nil,
	).WithLoadMode(goanalysis.LoadModeTypesInfo)
}
