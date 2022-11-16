package golinters

import (
	"github.com/Deng-Xian-Sheng/gtimeFormatLinter/pkg/analyzer"
	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"
	"golang.org/x/tools/go/analysis"
)

func NewGtimeFormatLinter() *goanalysis.Linter {
	return goanalysis.NewLinter(
		"gtimeformatlinter",
		"gtime.Time is the time type of the Go Frame framework. The formal parameters of its Format method are completely different from those of the Format method in the standard library.",
		[]*analysis.Analyzer{analyzer.Analyzer},
		nil,
	).WithLoadMode(goanalysis.LoadModeSyntax)
}
