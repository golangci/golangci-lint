package constructorcheck

import (
	"github.com/golangci/golangci-lint/pkg/goanalysis"
	constructorcheck "github.com/reflechant/constructor-check/analyzer"
	"golang.org/x/tools/go/analysis"
)

func New() *goanalysis.Linter {
	a := constructorcheck.Analyzer

	return goanalysis.NewLinter(
		a.Name,
		a.Doc,
		[]*analysis.Analyzer{a},
		nil,
	).WithLoadMode(goanalysis.LoadModeSyntax)
	// WithLoadMode(goanalysis.LoadModeTypesInfo)
}
