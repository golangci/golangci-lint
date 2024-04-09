package makezero

import (
	"fmt"
	"sync"

	"github.com/ashanbrown/makezero/makezero"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/goanalysis"
	"github.com/golangci/golangci-lint/pkg/lint/linter"
	"github.com/golangci/golangci-lint/pkg/result"
)

const name = "makezero"

func New(settings *config.MakezeroSettings) *goanalysis.Linter {
	var mu sync.Mutex
	var resIssues []goanalysis.Issue

	analyzer := &analysis.Analyzer{
		Name: name,
		Doc:  goanalysis.TheOnlyanalyzerDoc,
		Run: func(pass *analysis.Pass) (any, error) {
			issues, err := runMakeZero(pass, settings)
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
		name,
		"Finds slice declarations with non-zero initial length",
		[]*analysis.Analyzer{analyzer},
		nil,
	).WithIssuesReporter(func(*linter.Context) []goanalysis.Issue {
		return resIssues
	}).WithLoadMode(goanalysis.LoadModeTypesInfo)
}

func runMakeZero(pass *analysis.Pass, settings *config.MakezeroSettings) ([]goanalysis.Issue, error) {
	zero := makezero.NewLinter(settings.Always)

	var issues []goanalysis.Issue

	for _, file := range pass.Files {
		hints, err := zero.Run(pass.Fset, pass.TypesInfo, file)
		if err != nil {
			return nil, fmt.Errorf("makezero linter failed on file %q: %w", file.Name.String(), err)
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
