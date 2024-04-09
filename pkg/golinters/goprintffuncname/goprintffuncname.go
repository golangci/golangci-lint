package goprintffuncname

import (
	"github.com/jirfag/go-printf-func-name/pkg/analyzer"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/goanalysis"
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
