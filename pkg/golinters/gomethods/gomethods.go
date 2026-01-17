package gomethods

import (
	"github.com/golangci/golangci-lint/v2/pkg/goanalysis"

	"go.ufukty.com/gomethods/pkg/analyzer"
)

func New() *goanalysis.Linter {
	return goanalysis.NewLinterFromAnalyzer(analyzer.New())
}
