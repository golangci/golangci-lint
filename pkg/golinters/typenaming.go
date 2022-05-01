package golinters

import (
	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"
	typenaming "github.com/typenaming/typenaming/pkg/analyzer"
	"golang.org/x/tools/go/analysis"
)

func NewTypenaming() *goanalysis.Linter {
	return goanalysis.NewLinter(
		typenaming.Analyzer.Name,
		typenaming.Analyzer.Doc,
		[]*analysis.Analyzer{typenaming.Analyzer},
		nil,
	)
}
