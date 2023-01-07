package golinters

import (
	"github.com/atzoum/reuseconn/reuseconn"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"
)

func NewReuseconn() *goanalysis.Linter {
	return goanalysis.NewLinter(
		"reuseconn",
		"checks whether HTTP response body is consumed and closed properly in a single function",
		[]*analysis.Analyzer{reuseconn.Analyzer},
		nil,
	).WithLoadMode(goanalysis.LoadModeWholeProgram)
}
