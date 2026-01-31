package golistics

import (
	"github.com/golangci/golangci-lint/v2/pkg/goanalysis"

	"go.ufukty.com/golistics/pkg/analyzer"
)

func New() *goanalysis.Linter {
	return goanalysis.NewLinterFromAnalyzer(analyzer.New())
}
