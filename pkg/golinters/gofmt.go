package golinters

import (
	"sync"

	gofmtAPI "github.com/golangci/gofmt/gofmt"
	"github.com/pkg/errors"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"
	"github.com/golangci/golangci-lint/pkg/lint/linter"
)

const gofmtName = "gofmt"

func NewGofmt(settings *config.GoFmtSettings) *goanalysis.Linter {
	var mu sync.Mutex
	var resIssues []goanalysis.Issue

	analyzer := &analysis.Analyzer{
		Name: gofmtName,
		Doc:  goanalysis.TheOnlyanalyzerDoc,
		Run:  goanalysis.DummyRun,
	}

	return goanalysis.NewLinter(
		gofmtName,
		"Gofmt checks whether code was gofmt-ed. By default "+
			"this tool runs with -s option to check for code simplification",
		[]*analysis.Analyzer{analyzer},
		nil,
	).WithContextSetter(func(lintCtx *linter.Context) {
		analyzer.Run = func(pass *analysis.Pass) (interface{}, error) {
			issues, err := runGofmt(lintCtx, pass, settings)
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
		}
	}).WithIssuesReporter(func(*linter.Context) []goanalysis.Issue {
		return resIssues
	}).WithLoadMode(goanalysis.LoadModeSyntax)
}

func runGofmt(lintCtx *linter.Context, pass *analysis.Pass, settings *config.GoFmtSettings) ([]goanalysis.Issue, error) {
	fileNames := getFileNames(pass)

	var rewriteRules []gofmtAPI.RewriteRule
	for _, rule := range settings.RewriteRules {
		rewriteRules = append(rewriteRules, gofmtAPI.RewriteRule(rule))
	}

	var issues []goanalysis.Issue

	for _, f := range fileNames {
		diff, err := gofmtAPI.RunRewrite(f, settings.Simplify, rewriteRules)
		if err != nil { // TODO: skip
			return nil, err
		}
		if diff == nil {
			continue
		}

		is, err := extractIssuesFromPatch(string(diff), lintCtx, gofmtName)
		if err != nil {
			return nil, errors.Wrapf(err, "can't extract issues from gofmt diff output %q", string(diff))
		}

		for i := range is {
			issues = append(issues, goanalysis.NewIssue(&is[i], pass))
		}
	}

	return issues, nil
}
