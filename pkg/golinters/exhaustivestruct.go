package golinters

import (
	"github.com/mbilski/exhaustivestruct/pkg/analyzer"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"
)

func NewExhaustiveStruct() *goanalysis.Linter {
	return goanalysis.NewLinter(
		"exhaustivestruct",
		"Checks if all struct's fields are initialized",
		[]*analysis.Analyzer{analyzer.Analyzer},
		nil,
	).WithLoadMode(goanalysis.LoadModeTypesInfo)
}
