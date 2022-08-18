package golinters

import (
	"fmt"
	"sync"

	"golang.org/x/tools/go/analysis"
	"honnef.co/go/tools/unused"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"
	"github.com/golangci/golangci-lint/pkg/lint/linter"
	"github.com/golangci/golangci-lint/pkg/result"
)

const unusedName = "unused"

type UnusedSettings struct {
	GoVersion string
}

func NewUnused(settings *config.StaticCheckSettings) *goanalysis.Linter {
	var mu sync.Mutex
	var resIssues []goanalysis.Issue

	analyzer := &analysis.Analyzer{
		Name:     unusedName,
		Doc:      unused.Analyzer.Analyzer.Doc,
		Requires: unused.Analyzer.Analyzer.Requires,
		Run: func(pass *analysis.Pass) (interface{}, error) {
			issues, err := runUnused(pass)
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

	setAnalyzerGoVersion(analyzer, getGoVersion(settings))

	return goanalysis.NewLinter(
		unusedName,
		"Checks Go code for unused constants, variables, functions and types",
		[]*analysis.Analyzer{analyzer},
		nil,
	).WithIssuesReporter(func(lintCtx *linter.Context) []goanalysis.Issue {
		return resIssues
	}).WithLoadMode(goanalysis.LoadModeTypesInfo)
}

func runUnused(pass *analysis.Pass) ([]goanalysis.Issue, error) {
	res, err := unused.Analyzer.Analyzer.Run(pass)
	if err != nil {
		return nil, err
	}

	sr := unused.Serialize(pass, res.(unused.Result), pass.Fset)

	used := make(map[string]bool)
	for _, obj := range sr.Used {
		used[fmt.Sprintf("%s %d %s", obj.Position.Filename, obj.Position.Line, obj.Name)] = true
	}

	var issues []goanalysis.Issue

	// Inspired by https://github.com/dominikh/go-tools/blob/d694aadcb1f50c2d8ac0a1dd06217ebb9f654764/lintcmd/lint.go#L177-L197
	for _, object := range sr.Unused {
		if object.Kind == "type param" {
			continue
		}

		if object.InGenerated {
			continue
		}

		key := fmt.Sprintf("%s %d %s", object.Position.Filename, object.Position.Line, object.Name)
		if used[key] {
			continue
		}

		issue := goanalysis.NewIssue(&result.Issue{
			FromLinter: unusedName,
			Text:       fmt.Sprintf("%s %s is unused", object.Kind, object.Name),
			Pos:        object.Position,
		}, pass)

		issues = append(issues, issue)
	}

	return issues, nil
}
