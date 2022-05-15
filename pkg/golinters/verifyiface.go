package golinters

import (
	"github.com/dcu/verifyiface/analyzer"
	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"
	"golang.org/x/tools/go/analysis"
)

func NewVerifyIface() *goanalysis.Linter {
	linter := goanalysis.NewLinter(
		"verifyiface",
		"check that a interface implementation is verified according as explained here: https://github.com/uber-go/guide/blob/master/style.md#verify-interface-compliance",
		[]*analysis.Analyzer{analyzer.Analyzer},
		nil,
	).WithLoadMode(goanalysis.LoadModeTypesInfo)

	return linter
}
