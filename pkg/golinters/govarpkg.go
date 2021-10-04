package golinters

import (
	"github.com/alexal/govarpkg/pkg/analyzer"
	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"
	"golang.org/x/tools/go/analysis"
)

func NewGoVarPkg() *goanalysis.Linter {
	return goanalysis.NewLinter(
		"govarpkg",
		"Finds variables that collide with imported package names.",
		[]*analysis.Analyzer{analyzer.Analyzer},
		nil,
	).WithLoadMode(goanalysis.LoadModeSyntax)
}
