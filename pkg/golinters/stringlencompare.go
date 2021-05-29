package golinters

import (
	"github.com/johejo/stringlencompare"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"
)

func NewStringLenCompare() *goanalysis.Linter {
	return goanalysis.NewLinter(
		"stringlencompare",
		"check string len compare style",
		[]*analysis.Analyzer{stringlencompare.Analyzer},
		nil,
	).WithLoadMode(goanalysis.LoadModeNone)
}
