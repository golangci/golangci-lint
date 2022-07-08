package golinters

import (
	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"

	"github.com/timakin/bodyclose/passes/bodyclose"
	"golang.org/x/tools/go/analysis"
)

func NewBodyclose() *goanalysis.Linter {
	return goanalysis.NewLinter(
		"bodyclose",
		"checks whether HTTP response body is closed successfully",
		[]*analysis.Analyzer{bodyclose.Analyzer},
		nil,
	).WithLoadMode(goanalysis.LoadModeTypesInfo)
}
