package golinters

import (
	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"

	"github.com/sivchari/containedctx"
	"golang.org/x/tools/go/analysis"
)

func NewContainedCtx() *goanalysis.Linter {
	a := containedctx.Analyzer

	return goanalysis.NewLinter(
		a.Name,
		a.Doc,
		[]*analysis.Analyzer{a},
		nil,
	).WithLoadMode(goanalysis.LoadModeSyntax)
}
