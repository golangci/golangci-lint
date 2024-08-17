package scancheck

import (
	"github.com/golangci/golangci-lint/pkg/goanalysis"
	"github.com/raidancampbell/scancheck/pkg/scancheck"
	"golang.org/x/tools/go/analysis"
)

func New() *goanalysis.Linter {
	a := scancheck.Analyzer

	return goanalysis.NewLinter(
		a.Name,
		a.Doc,
		[]*analysis.Analyzer{a},
		nil,
	).WithLoadMode(goanalysis.LoadModeTypesInfo)
}
