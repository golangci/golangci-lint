package golinters

import (
	"honnef.co/go/tools/stylecheck"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"
)

func NewStylecheck(settings *config.StaticCheckSettings) *goanalysis.Linter {
	analyzers := setupStaticCheckAnalyzers(stylecheck.Analyzers, settings)

	return goanalysis.NewLinter(
		"stylecheck",
		"Stylecheck is a replacement for golint",
		analyzers,
		nil,
	).WithLoadMode(goanalysis.LoadModeTypesInfo)
}
