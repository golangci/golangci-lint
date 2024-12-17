package mirror

import (
	"github.com/butuzov/mirror"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/goanalysis"
)

func New() *goanalysis.Linter {
	a := mirror.NewAnalyzer()

	// mirror only lints test files if the `--with-tests` flag is passed,
	// so we pass the `with-tests` flag as true to the analyzer before running it.
	// This can be turned off by using the regular golangci-lint flags such as `--tests` or `--skip-files`
	// or can be disabled per linter via exclude rules.
	// (see https://github.com/golangci/golangci-lint/issues/2527#issuecomment-1023707262)
	linterCfg := map[string]map[string]any{
		a.Name: {
			"with-tests": true,
		},
	}

	return goanalysis.NewLinter(
		a.Name,
		a.Doc,
		[]*analysis.Analyzer{a},
		linterCfg,
	).WithLoadMode(goanalysis.LoadModeTypesInfo)
}
