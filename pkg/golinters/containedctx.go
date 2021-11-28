package golinters

import (
	"github.com/sivchari/containedctx"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"
)

func NewContainedCtx() *goanalysis.Linter {
	a := containedctx.Analyzer

	analyzers := []*analysis.Analyzer{
		a,
	}

	return goanalysis.NewLinter(
		a.Name,
		a.Doc,
		analyzers,
		nil,
	).WithLoadMode(goanalysis.LoadModeSyntax)
}
