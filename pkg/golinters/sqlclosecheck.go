package golinters

import (
	"github.com/ryanrolds/sqlclosecheck/pkg/analyzer"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"
)

func NewSQLCloseCheck() *goanalysis.Linter {
	return goanalysis.NewLinter(
		"sqlclosecheck",
		"Checks that sql.Rows and sql.Stmt are closed.",
		[]*analysis.Analyzer{
			analyzer.NewAnalyzer(),
		},
		nil,
	).WithLoadMode(goanalysis.LoadModeTypesInfo)
}
