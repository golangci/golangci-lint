package whyvarscope

import (
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/goanalysis"
	"github.com/itsamirhn/whyvarscope/pkg/analyzer"
)

func New() *goanalysis.Linter {
	a := analyzer.Analyzer

	return goanalysis.NewLinter(
		a.Name,
		a.Doc,
		[]*analysis.Analyzer{a},
		nil,
	).WithLoadMode(goanalysis.LoadModeSyntax)
}
