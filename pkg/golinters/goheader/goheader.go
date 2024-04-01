package goheader

import (
	"go/token"
	"sync"

	goheader "github.com/denis-tingaikin/go-header"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/goanalysis"
	"github.com/golangci/golangci-lint/pkg/lint/linter"
	"github.com/golangci/golangci-lint/pkg/result"
)

const name = "goheader"

func New(settings *config.GoHeaderSettings) *goanalysis.Linter {
	var mu sync.Mutex
	var resIssues []goanalysis.Issue

	conf := &goheader.Configuration{}
	if settings != nil {
		conf = &goheader.Configuration{
			Values:       settings.Values,
			Template:     settings.Template,
			TemplatePath: settings.TemplatePath,
		}
	}

	analyzer := &analysis.Analyzer{
		Name: name,
		Doc:  goanalysis.TheOnlyanalyzerDoc,
		Run: func(pass *analysis.Pass) (any, error) {
			issues, err := runGoHeader(pass, conf)
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
		"Checks is file header matches to pattern",
		[]*analysis.Analyzer{analyzer},
		nil,
	).WithIssuesReporter(func(*linter.Context) []goanalysis.Issue {
		return resIssues
	}).WithLoadMode(goanalysis.LoadModeSyntax)
}

func runGoHeader(pass *analysis.Pass, conf *goheader.Configuration) ([]goanalysis.Issue, error) {
	if conf.TemplatePath == "" && conf.Template == "" {
		// User did not pass template, so then do not run go-header linter
		return nil, nil
	}

	template, err := conf.GetTemplate()
	if err != nil {
		return nil, err
	}

	values, err := conf.GetValues()
	if err != nil {
		return nil, err
	}

	a := goheader.New(goheader.WithTemplate(template), goheader.WithValues(values))

	var issues []goanalysis.Issue
	for _, file := range pass.Files {
		path := pass.Fset.Position(file.Pos()).Filename

		i := a.Analyze(&goheader.Target{File: file, Path: path})

		if i == nil {
			continue
		}

		issue := result.Issue{
			Pos: token.Position{
				Line:     i.Location().Line + 1,
				Column:   i.Location().Position,
				Filename: path,
			},
			Text:       i.Message(),
			FromLinter: name,
		}

		if fix := i.Fix(); fix != nil {
			issue.LineRange = &result.Range{
				From: issue.Line(),
				To:   issue.Line() + len(fix.Actual) - 1,
			}
			issue.Replacement = &result.Replacement{
				NeedOnlyDelete: len(fix.Expected) == 0,
				NewLines:       fix.Expected,
			}
		}

		issues = append(issues, goanalysis.NewIssue(&issue, pass))
	}

	return issues, nil
}
