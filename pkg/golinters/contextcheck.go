package golinters

import (
	"github.com/sylvia7788/contextcheck"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"
)

func NewContextCheck() *goanalysis.Linter {
	analyzer := contextcheck.NewAnalyzer()
	return goanalysis.NewLinter(
		"contextcheck",
		"check for using context.Background() and context.TODO() directly",
		[]*analysis.Analyzer{analyzer},
		nil,
	).WithLoadMode(goanalysis.LoadModeTypesInfo)
}
