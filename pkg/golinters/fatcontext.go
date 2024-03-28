package golinters

import (
	"github.com/Crocmagnon/fatcontext/pkg/analyzer"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/goanalysis"
)

func NewFatContext() *goanalysis.Linter {
	return goanalysis.NewLinter(
		"fatcontext",
		"Detects potential fat contexts in loops",
		[]*analysis.Analyzer{analyzer.Analyzer},
		nil,
	).WithLoadMode(goanalysis.LoadModeTypesInfo)
}
