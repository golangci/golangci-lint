package golinters

import (
	"github.com/sivchari/copywaitgroup"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"
)

func NewCopyWaitGroup() *goanalysis.Linter {
	analyzers := []*analysis.Analyzer{
		copywaitgroup.Analyzer,
	}

	return goanalysis.NewLinter(
		"copywaitgroup",
		"finds a func that passes sync.WaitGroup as a value instead of a pointer",
		analyzers,
		nil,
	).WithLoadMode(goanalysis.LoadModeTypesInfo)
}
