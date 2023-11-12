package golinters

import (
	"github.com/hendrywiranto/gomockcontrollerfinish/pkg/analyzer"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"
)

func NewGomockControllerFinish() *goanalysis.Linter {
	analyzers := []*analysis.Analyzer{
		analyzer.New(),
	}

	return goanalysis.NewLinter(
		"gomockcontrollerfinish",
		"Checks whether an unnecessary call to .Finish() on gomock.Controller exists",
		analyzers,
		nil,
	).WithLoadMode(goanalysis.LoadModeTypesInfo)
}
