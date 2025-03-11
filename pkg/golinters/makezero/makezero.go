package makezero

import (
	"fmt"

	"github.com/ashanbrown/makezero/makezero"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/v2/pkg/config"
	"github.com/golangci/golangci-lint/v2/pkg/goanalysis"
)

const linterName = "makezero"

func New(settings *config.MakezeroSettings) *goanalysis.Linter {
	analyzer := &analysis.Analyzer{
		Name: linterName,
		Doc:  goanalysis.TheOnlyanalyzerDoc,
		Run: func(pass *analysis.Pass) (any, error) {
			err := runMakeZero(pass, settings)
			if err != nil {
				return nil, err
			}

			return nil, nil
		},
	}

	return goanalysis.NewLinter(
		linterName,
		"Finds slice declarations with non-zero initial length",
		[]*analysis.Analyzer{analyzer},
		nil,
	).WithLoadMode(goanalysis.LoadModeTypesInfo)
}

func runMakeZero(pass *analysis.Pass, settings *config.MakezeroSettings) error {
	zero := makezero.NewLinter(settings.Always)

	for _, file := range pass.Files {
		hints, err := zero.Run(pass.Fset, pass.TypesInfo, file)
		if err != nil {
			return fmt.Errorf("makezero linter failed on file %q: %w", file.Name.String(), err)
		}

		for _, hint := range hints {
			pass.Report(analysis.Diagnostic{
				Pos:     hint.Pos(),
				Message: hint.Details(),
			})
		}
	}

	return nil
}
