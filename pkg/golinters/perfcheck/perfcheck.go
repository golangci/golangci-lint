package perfcheck

import (
	"golang.org/x/tools/go/analysis"

	"github.com/m-v-kalashnikov/perfcheck/go/pkg/perfchecklint"

	"github.com/golangci/golangci-lint/v2/pkg/goanalysis"
)

func New() *goanalysis.Linter {
	linter := perfchecklint.BuildGoanalysis(newGoanalysisLinter, perfchecklint.BuildOptions{})

	return linter.WithLoadMode(goanalysis.LoadModeTypesInfo)
}

func newGoanalysisLinter(name, desc string, analyzers ...*analysis.Analyzer) *goanalysis.Linter {
	return goanalysis.NewLinter(name, desc, analyzers, nil)
}
