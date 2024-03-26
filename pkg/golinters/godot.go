package golinters

import (
	"sync"

	"github.com/tetafro/godot"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/goanalysis"
	"github.com/golangci/golangci-lint/pkg/lint/linter"
	"github.com/golangci/golangci-lint/pkg/result"
)

const godotName = "godot"

func NewGodot(settings *config.GodotSettings) *goanalysis.Linter {
	var mu sync.Mutex
	var resIssues []goanalysis.Issue

	var dotSettings godot.Settings

	if settings != nil {
		dotSettings = godot.Settings{
			Scope:   godot.Scope(settings.Scope),
			Exclude: settings.Exclude,
			Period:  settings.Period,
			Capital: settings.Capital,
		}

		// Convert deprecated setting
		if settings.CheckAll {
			dotSettings.Scope = godot.AllScope
		}
	}

	if dotSettings.Scope == "" {
		dotSettings.Scope = godot.DeclScope
	}

	analyzer := &analysis.Analyzer{
		Name: godotName,
		Doc:  goanalysis.TheOnlyanalyzerDoc,
		Run: func(pass *analysis.Pass) (any, error) {
			issues, err := runGodot(pass, dotSettings)
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

	return goanalysis.NewLinter(
		godotName,
		"Check if comments end in a period",
		[]*analysis.Analyzer{analyzer},
		nil,
	).WithIssuesReporter(func(*linter.Context) []goanalysis.Issue {
		return resIssues
	}).WithLoadMode(goanalysis.LoadModeSyntax)
}

func runGodot(pass *analysis.Pass, settings godot.Settings) ([]goanalysis.Issue, error) {
	var lintIssues []godot.Issue
	for _, file := range pass.Files {
		iss, err := godot.Run(file, pass.Fset, settings)
		if err != nil {
			return nil, err
		}
		lintIssues = append(lintIssues, iss...)
	}

	if len(lintIssues) == 0 {
		return nil, nil
	}

	issues := make([]goanalysis.Issue, len(lintIssues))
	for k, i := range lintIssues {
		issue := result.Issue{
			Pos:        i.Pos,
			Text:       i.Message,
			FromLinter: godotName,
			Replacement: &result.Replacement{
				NewLines: []string{i.Replacement},
			},
		}

		issues[k] = goanalysis.NewIssue(&issue, pass)
	}

	return issues, nil
}
