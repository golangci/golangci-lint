package forbidigo

import (
	"fmt"
	"sync"

	"github.com/ashanbrown/forbidigo/forbidigo"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/goanalysis"
	"github.com/golangci/golangci-lint/pkg/lint/linter"
	"github.com/golangci/golangci-lint/pkg/logutils"
	"github.com/golangci/golangci-lint/pkg/result"
)

const name = "forbidigo"

func New(settings *config.ForbidigoSettings) *goanalysis.Linter {
	var mu sync.Mutex
	var resIssues []goanalysis.Issue

	analyzer := &analysis.Analyzer{
		Name: name,
		Doc:  goanalysis.TheOnlyanalyzerDoc,
		Run: func(pass *analysis.Pass) (any, error) {
			issues, err := runForbidigo(pass, settings)
			if err != nil {
				return nil, err
			}

			if len(issues) == 0 {
				return nil, nil
			}

			mu.Lock()
			resIssues = append(resIssues, issues...)
			mu.Unlock()
			return nil, nil
		},
	}

	// Without AnalyzeTypes, LoadModeSyntax is enough.
	// But we cannot make this depend on the settings and have to mirror the mode chosen in GetAllSupportedLinterConfigs,
	// therefore we have to use LoadModeTypesInfo in all cases.
	return goanalysis.NewLinter(
		name,
		"Forbids identifiers",
		[]*analysis.Analyzer{analyzer},
		nil,
	).WithIssuesReporter(func(*linter.Context) []goanalysis.Issue {
		return resIssues
	}).WithLoadMode(goanalysis.LoadModeTypesInfo)
}

func runForbidigo(pass *analysis.Pass, settings *config.ForbidigoSettings) ([]goanalysis.Issue, error) {
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
			return nil, err
		}
		patterns = append(patterns, string(buffer))
	}

	forbid, err := forbidigo.NewLinter(patterns, options...)
	if err != nil {
		return nil, fmt.Errorf("failed to create linter %q: %w", name, err)
	}

	var issues []goanalysis.Issue
	for _, file := range pass.Files {
		runConfig := forbidigo.RunConfig{
			Fset:     pass.Fset,
			DebugLog: logutils.Debug(logutils.DebugKeyForbidigo),
		}
		if settings != nil && settings.AnalyzeTypes {
			runConfig.TypesInfo = pass.TypesInfo
		}
		hints, err := forbid.RunWithConfig(runConfig, file)
		if err != nil {
			return nil, fmt.Errorf("forbidigo linter failed on file %q: %w", file.Name.String(), err)
		}

		for _, hint := range hints {
			issues = append(issues, goanalysis.NewIssue(&result.Issue{
				Pos:        hint.Position(),
				Text:       hint.Details(),
				FromLinter: name,
			}, pass))
		}
	}

	return issues, nil
}
