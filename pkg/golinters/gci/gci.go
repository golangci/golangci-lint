package gci

import (
	"bytes"
	"fmt"
	"io"
	"os"

	gcicfg "github.com/daixiang0/gci/pkg/config"
	"github.com/daixiang0/gci/pkg/gci"
	"github.com/daixiang0/gci/pkg/log"
	"github.com/shazow/go-diff/difflib"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/goanalysis"
	"github.com/golangci/golangci-lint/pkg/golinters/internal"
	"github.com/golangci/golangci-lint/pkg/lint/linter"
)

const linterName = "gci"

type differ interface {
	Diff(out io.Writer, a io.ReadSeeker, b io.ReadSeeker) error
}

func New(settings *config.GciSettings) *goanalysis.Linter {
	log.InitLogger()
	_ = log.L().Sync()

	diff := difflib.New()

	a := &analysis.Analyzer{
		Name: linterName,
		Doc:  goanalysis.TheOnlyanalyzerDoc,
		Run:  goanalysis.DummyRun,
	}

	return goanalysis.NewLinter(
		linterName,
		"Checks if code and import statements are formatted, it makes import statements always deterministic.",
		[]*analysis.Analyzer{a},
		nil,
	).WithContextSetter(func(lintCtx *linter.Context) {
		a.Run = func(pass *analysis.Pass) (any, error) {
			err := run(lintCtx, pass, settings, diff)
			if err != nil {
				return nil, err
			}

			return nil, nil
		}
	}).WithLoadMode(goanalysis.LoadModeSyntax)
}

func run(lintCtx *linter.Context, pass *analysis.Pass, settings *config.GciSettings, diff differ) error {
	cfg := gcicfg.YamlConfig{
		Cfg: gcicfg.BoolConfig{
			NoInlineComments: settings.NoInlineComments,
			NoPrefixComments: settings.NoPrefixComments,
			SkipGenerated:    settings.SkipGenerated,
			CustomOrder:      settings.CustomOrder,
			NoLexOrder:       settings.NoLexOrder,
		},
		SectionStrings: settings.Sections,
		ModPath:        pass.Module.Path,
	}

	if settings.LocalPrefixes != "" {
		cfg.SectionStrings = []string{
			"standard",
			"default",
			fmt.Sprintf("prefix(%s)", settings.LocalPrefixes),
		}
	}

	parsedCfg, err := cfg.Parse()
	if err != nil {
		return err
	}

	for _, file := range pass.Files {
		position, isGoFile := goanalysis.GetGoFilePosition(pass, file)
		if !isGoFile {
			continue
		}

		input, err := os.ReadFile(position.Filename)
		if err != nil {
			return fmt.Errorf("unable to open file %s: %w", position.Filename, err)
		}

		_, output, err := gci.LoadFormat(input, position.Filename, *parsedCfg)
		if err != nil {
			return fmt.Errorf("error while running gci: %w", err)
		}

		if !bytes.Equal(input, output) {
			out := bytes.NewBufferString(fmt.Sprintf("--- %[1]s\n+++ %[1]s\n", position.Filename))

			err := diff.Diff(out, bytes.NewReader(input), bytes.NewReader(output))
			if err != nil {
				return fmt.Errorf("error while running gci: %w", err)
			}

			diff := out.String()

			err = internal.ExtractDiagnosticFromPatch(pass, file, diff, lintCtx)
			if err != nil {
				return fmt.Errorf("can't extract issues from gci diff output %q: %w", diff, err)
			}
		}
	}

	return nil
}
