package golinters

import (
	"github.com/lasiar/canonicalheader"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"
)

func NewCanonicalheder() *goanalysis.Linter {
	a := canonicalheader.Analyzer

	return goanalysis.NewLinter(
		a.Name,
		a.Doc,
		[]*analysis.Analyzer{canonicalheader.Analyzer},
		nil,
	).WithLoadMode(goanalysis.LoadModeTypesInfo)
}
