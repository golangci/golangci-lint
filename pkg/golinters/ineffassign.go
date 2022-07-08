package golinters

import (
	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"

	"github.com/gordonklaus/ineffassign/pkg/ineffassign"
	"golang.org/x/tools/go/analysis"
)

func NewIneffassign() *goanalysis.Linter {
	return goanalysis.NewLinter(
		"ineffassign",
		"Detects when assignments to existing variables are not used",
		[]*analysis.Analyzer{ineffassign.Analyzer},
		nil,
	).WithLoadMode(goanalysis.LoadModeSyntax)
}
