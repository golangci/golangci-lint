package golinters

import (
	"github.com/dorfire/go-analyzers/src/visibilityorder"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"
)

func NewVisibilityOrder() *goanalysis.Linter {
	return goanalysis.NewLinter(
		visibilityorder.Analyzer.Name,
		visibilityorder.Analyzer.Doc,
		[]*analysis.Analyzer{visibilityorder.Analyzer},
		nil,
	).WithLoadMode(goanalysis.LoadModeSyntax)
}
