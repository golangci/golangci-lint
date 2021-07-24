package golinters

import (
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"
	"github.com/sivchari/nilassign"
)

func NewNilAssign() *goanalysis.Linter {
	analyzers := []*analysis.Analyzer{
		nilassign.Analyzer,
	}

	return goanalysis.NewLinter(
		"nilassign",
		"Finds that assigning to invalid memory address or nil pointer dereference.",
		analyzers,
		nil,
	).WithLoadMode(goanalysis.LoadModeSyntax)
}
