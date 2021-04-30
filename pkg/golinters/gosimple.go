package golinters

import (
	"honnef.co/go/tools/simple"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"
)

func NewGosimple(settings *config.StaticCheckSettings) *goanalysis.Linter {
	analyzers := setupStaticCheckAnalyzers(simple.Analyzers, settings)

	return goanalysis.NewLinter(
		"gosimple",
		"Linter for Go source code that specializes in simplifying a code",
		analyzers,
		nil,
	).WithLoadMode(goanalysis.LoadModeTypesInfo)
}
