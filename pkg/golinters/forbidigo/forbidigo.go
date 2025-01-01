package forbidigo

import (
	"fmt"

	"github.com/ashanbrown/forbidigo/forbidigo"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/goanalysis"
	"github.com/golangci/golangci-lint/pkg/logutils"
)

const linterName = "forbidigo"

func New(settings *config.ForbidigoSettings) *goanalysis.Linter {
	analyzer := &analysis.Analyzer{
		Name: linterName,
		Doc:  goanalysis.TheOnlyanalyzerDoc,
		Run: func(pass *analysis.Pass) (any, error) {
			err := runForbidigo(pass, settings)
			if err != nil {
				return nil, err
			}

			return nil, nil
		},
	}

	// Without AnalyzeTypes, LoadModeSyntax is enough.
	// But we cannot make this depend on the settings and have to mirror the mode chosen in GetAllSupportedLinterConfigs,
	// therefore we have to use LoadModeTypesInfo in all cases.
	return goanalysis.NewLinter(
		linterName,
		"Forbids identifiers",
		[]*analysis.Analyzer{analyzer},
		nil,
	).WithLoadMode(goanalysis.LoadModeTypesInfo)
}

func runForbidigo(pass *analysis.Pass, settings *config.ForbidigoSettings) error {
	options := []forbidigo.Option{
		forbidigo.OptionExcludeGodocExamples(settings.ExcludeGodocExamples),
		// disable "//permit" directives so only "//nolint" directives matters within golangci-lint
		forbidigo.OptionIgnorePermitDirectives(true),
		forbidigo.OptionAnalyzeTypes(settings.AnalyzeTypes),
	}

	// Convert patterns back to strings because that is what NewLinter accepts.
	var patterns []string
	for _, pattern := range settings.Forbid {
		buffer, err := pattern.MarshalString()
		if err != nil {
			return err
		}

		patterns = append(patterns, string(buffer))
	}

	forbid, err := forbidigo.NewLinter(patterns, options...)
	if err != nil {
		return fmt.Errorf("failed to create linter %q: %w", linterName, err)
	}

	for _, file := range pass.Files {
		runConfig := forbidigo.RunConfig{
			Fset:     pass.Fset,
			DebugLog: logutils.Debug(logutils.DebugKeyForbidigo),
		}

		if settings.AnalyzeTypes {
			runConfig.TypesInfo = pass.TypesInfo
		}

		hints, err := forbid.RunWithConfig(runConfig, file)
		if err != nil {
			return fmt.Errorf("forbidigo linter failed on file %q: %w", file.Name.String(), err)
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
