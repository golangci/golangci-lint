package goconst

import (
	"fmt"
	"sync"

	goconstAPI "github.com/jgautheron/goconst"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/v2/pkg/config"
	"github.com/golangci/golangci-lint/v2/pkg/goanalysis"
	"github.com/golangci/golangci-lint/v2/pkg/golinters/internal"
	"github.com/golangci/golangci-lint/v2/pkg/lint/linter"
	"github.com/golangci/golangci-lint/v2/pkg/result"
)

const linterName = "goconst"

func New(settings *config.GoConstSettings) *goanalysis.Linter {
	var mu sync.Mutex
	var resIssues []goanalysis.Issue

	analyzer := &analysis.Analyzer{
		Name: linterName,
		Doc:  goanalysis.TheOnlyanalyzerDoc,
		Run: func(pass *analysis.Pass) (any, error) {
			issues, err := runGoconst(pass, settings)
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
		linterName,
		"Finds repeated strings that could be replaced by a constant",
		[]*analysis.Analyzer{analyzer},
		nil,
	).WithIssuesReporter(func(*linter.Context) []goanalysis.Issue {
		return resIssues
	}).WithLoadMode(goanalysis.LoadModeSyntax)
}

func runGoconst(pass *analysis.Pass, settings *config.GoConstSettings) ([]goanalysis.Issue, error) {
	cfg := goconstAPI.Config{
		IgnoreStrings:      settings.IgnoreStrings,
		MatchWithConstants: settings.MatchWithConstants,
		MinStringLength:    settings.MinStringLen,
		MinOccurrences:     settings.MinOccurrencesCount,
		ParseNumbers:       settings.ParseNumbers,
		NumberMin:          settings.NumberMin,
		NumberMax:          settings.NumberMax,
		ExcludeTypes:       map[goconstAPI.Type]bool{},

		// Should be managed with `linters.exclusions.rules`.
		IgnoreTests: false,
	}

	if settings.IgnoreCalls {
		cfg.ExcludeTypes[goconstAPI.Call] = true
	}

	lintIssues, err := goconstAPI.Run(pass.Files, pass.Fset, &cfg)
	if err != nil {
		return nil, err
	}

	if len(lintIssues) == 0 {
		return nil, nil
	}

	res := make([]goanalysis.Issue, 0, len(lintIssues))
	for _, i := range lintIssues {
		text := fmt.Sprintf("string %s has %d occurrences", internal.FormatCode(i.Str, nil), i.OccurrencesCount)

		if i.MatchingConst == "" {
			text += ", make it a constant"
		} else {
			text += fmt.Sprintf(", but such constant %s already exists", internal.FormatCode(i.MatchingConst, nil))
		}

		res = append(res, goanalysis.NewIssue(&result.Issue{
			Pos:        i.Pos,
			Text:       text,
			FromLinter: linterName,
		}, pass))
	}

	return res, nil
}
