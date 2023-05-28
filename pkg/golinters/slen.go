package golinters

import (
	"github.com/piotrpersona/slen/pkg/analyzer"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"
)

func NewSlen() *goanalysis.Linter {
	return goanalysis.NewLinter(
		analyzer.SlenCmd,
		analyzer.SlenDescription,
		[]*analysis.Analyzer{
			analyzer.Analyzer,
		},
		nil,
	)
}
