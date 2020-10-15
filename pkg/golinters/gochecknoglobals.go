package golinters

import (
	"golang.org/x/tools/go/analysis"

	"4d63.com/gochecknoglobals/checknoglobals"

	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"
)

func NewGochecknoglobals() *goanalysis.Linter {
	gochecknoglobals := checknoglobals.Analyzer()

	// gochecknoglobals only lints test files if the `-t` flag is passed so we
	// pass the `t` flag as true to the analyzer before running it. This can be
	// turned of by using the regular golangci-lint flags such as `--tests` or
	// `--skip-files`.
	linterConfig := map[string]map[string]interface{}{
		gochecknoglobals.Name: {
			"t": true,
		},
	}

	return goanalysis.NewLinter(
		gochecknoglobals.Name,
		gochecknoglobals.Doc,
		[]*analysis.Analyzer{gochecknoglobals},
		linterConfig,
	).WithLoadMode(goanalysis.LoadModeSyntax)
}
