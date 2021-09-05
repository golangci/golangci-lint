package golinters

import (
	"github.com/sivchari/nobreak"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"
)

func NewNobreak() *goanalysis.Linter {
	analyzers := []*analysis.Analyzer{
		nobreak.Analyzer,
	}

	return goanalysis.NewLinter(
		"nobreak",
		"Find inifinite `for statement` loop",
		analyzers,
		nil,
	).WithLoadMode(goanalysis.LoadModeSyntax)
}
