package golinters

import (
	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"
	"go.abhg.dev/unexportedglobal"
	"golang.org/x/tools/go/analysis"
)

func NewUnexportedGlobal() *goanalysis.Linter {
	return goanalysis.NewLinter(
		"unexportedglobal",
		"Disallows unexported globals without a '_' prefix",
		[]*analysis.Analyzer{unexportedglobal.Analyzer},
		nil,
	).WithLoadMode(goanalysis.LoadModeTypesInfo)
}
