package golinters

import (
	"github.com/macabu/inamedparam"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"
)

func NewINamedParam() *goanalysis.Linter {
	analyzer := inamedparam.Analyzer

	return goanalysis.NewLinter(
		analyzer.Name,
		analyzer.Doc,
		[]*analysis.Analyzer{
			analyzer,
		},
		nil,
	).WithLoadMode(goanalysis.LoadModeSyntax)
}
