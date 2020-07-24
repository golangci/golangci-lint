package golinters

import (
	"github.com/ssgreg/nlreturn/v2/pkg/nlreturn"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"
)

func NewNLReturn() *goanalysis.Linter {
	return goanalysis.NewLinter(
		"nlreturn",
		"nlreturn checks for a new line before return and branch statements to increase code clarity",
		[]*analysis.Analyzer{
			nlreturn.NewAnalyzer(),
		},
		nil,
	).WithLoadMode(goanalysis.LoadModeSyntax)
}
