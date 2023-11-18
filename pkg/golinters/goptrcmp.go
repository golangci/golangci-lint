package golinters

import (
	"github.com/w1ck3dg0ph3r/goptrcmp"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"
)

func NewGoPtrCmp() *goanalysis.Linter {
	return goanalysis.NewLinter(
		"goptrcmp",
		"Reports comparison between pointer values",
		[]*analysis.Analyzer{goptrcmp.Analyzer()},
		nil,
	).WithLoadMode(goanalysis.LoadModeTypesInfo)
}
