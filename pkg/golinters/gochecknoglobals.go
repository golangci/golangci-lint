package golinters

import (
	"flag"

	"golang.org/x/tools/go/analysis"

	"4d63.com/gochecknoglobals/checknoglobals"

	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"
)

func NewGochecknoglobals() *goanalysis.Linter {
	// gochecknoglobals only lints test files if the `-t` flag is passed so we
	// set up our own FlagSet and add it to the analyzer before running it. This
	// can be turned of by using the regular golangci-lint flags such as
	// `--tests` or `--skip-files`.
	flags := flag.NewFlagSet("", flag.ContinueOnError)
	flags.Bool("t", true, "Include tests")

	gochecknoglobals := checknoglobals.Analyzer()
	gochecknoglobals.Flags = *flags

	return goanalysis.NewLinter(
		gochecknoglobals.Name,
		gochecknoglobals.Doc,
		[]*analysis.Analyzer{gochecknoglobals},
		nil,
	).WithLoadMode(goanalysis.LoadModeSyntax)
}
