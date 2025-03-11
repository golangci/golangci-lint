package ineffassign

import (
	"github.com/gordonklaus/ineffassign/pkg/ineffassign"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/v2/pkg/goanalysis"
)

func New() *goanalysis.Linter {
	a := ineffassign.Analyzer

	return goanalysis.NewLinter(
		a.Name,
		"Detects when assignments to existing variables are not used",
		[]*analysis.Analyzer{a},
		nil,
	).WithLoadMode(goanalysis.LoadModeSyntax)
}
