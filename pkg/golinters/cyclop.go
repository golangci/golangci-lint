package golinters

import (
	"github.com/bkielbasa/cyclop/pkg/analyzer"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"
)

const cyclopName = "cyclop"

func NewCyclop() *goanalysis.Linter {
	return goanalysis.NewLinter(
		cyclopName,
		"checks function and package cyclomatic complexity",
		[]*analysis.Analyzer{
			analyzer.NewAnalyzer(),
		},
		nil,
	).WithLoadMode(goanalysis.LoadModeTypesInfo)
}
