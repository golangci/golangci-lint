package golinters

import (
	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"
	"github.com/le0tk0k/noioutil"
	"golang.org/x/tools/go/analysis"
)

func NewNoIoUtil() *goanalysis.Linter {
	analyzers := []*analysis.Analyzer{
		noioutil.Analyzer,
	}

	return goanalysis.NewLinter(
		"noioutil",
		"noioutil finds io/ioutil package",
		analyzers,
		nil,
	).WithLoadMode(goanalysis.LoadModeTypesInfo)
}
