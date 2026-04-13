package goconcurrencylint

import (
	concurrencyanalyzer "github.com/sanbricio/goconcurrencylint/pkg/analyzer"

	"github.com/golangci/golangci-lint/v2/pkg/goanalysis"
)

func New() *goanalysis.Linter {
	return goanalysis.
		NewLinterFromAnalyzer(concurrencyanalyzer.Analyzer).
		WithDesc("Detects incorrect sync.Mutex, sync.RWMutex, and sync.WaitGroup usage").
		WithLoadMode(goanalysis.LoadModeTypesInfo)
}
