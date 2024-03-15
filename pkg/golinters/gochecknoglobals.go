package golinters

import (
	"4d63.com/gochecknoglobals/checknoglobals"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"
)

func NewGochecknoglobals() *goanalysis.Linter {
	a := checknoglobals.Analyzer()

	// gochecknoglobals only lints test files if the `-t` flag is passed,
	// so we pass the `t` flag as true to the analyzer before running it.
	// This can be turned off by using the regular golangci-lint flags such as `--tests` or `--exclude-files`.
	linterConfig := map[string]map[string]any{
		a.Name: {"t": true},
	}

	return goanalysis.NewLinter(
		a.Name,
		"Check that no global variables exist.",
		[]*analysis.Analyzer{a},
		linterConfig,
	).WithLoadMode(goanalysis.LoadModeTypesInfo)
}
