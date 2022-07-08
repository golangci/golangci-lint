package golinters

import (
	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"

	"github.com/stbenjam/no-sprintf-host-port/pkg/analyzer"
	"golang.org/x/tools/go/analysis"
)

func NewNoSprintfHostPort() *goanalysis.Linter {
	a := analyzer.Analyzer

	return goanalysis.NewLinter(
		a.Name,
		a.Doc,
		[]*analysis.Analyzer{a},
		nil,
	).WithLoadMode(goanalysis.LoadModeSyntax)
}
