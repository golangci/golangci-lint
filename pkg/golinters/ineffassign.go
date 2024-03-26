package golinters

import (
	"github.com/gordonklaus/ineffassign/pkg/ineffassign"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/goanalysis"
)

func NewIneffassign() *goanalysis.Linter {
	a := ineffassign.Analyzer

	return goanalysis.NewLinter(
		a.Name,
		"Detects when assignments to existing variables are not used",
		[]*analysis.Analyzer{a},
		nil,
	).WithLoadMode(goanalysis.LoadModeSyntax)
}
