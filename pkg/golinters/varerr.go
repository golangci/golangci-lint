package golinters

import (
	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"
	"github.com/sivchari/varerr"
	"golang.org/x/tools/go/analysis"
)

func NewVarErr() *goanalysis.Linter {
	analyzers := []*analysis.Analyzer{
		varerr.Analyzer,
	}

	return goanalysis.NewLinter(
		"varerr",
		"Checks that you initialize error type.",
		analyzers,
		nil,
	).WithLoadMode(goanalysis.LoadModeSyntax)
}
