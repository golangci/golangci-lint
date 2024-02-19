package golinters

import (
	"fmt"
	"sync"

	"golang.org/x/tools/go/analysis"
	"honnef.co/go/tools/analysis/facts/directives"
	"honnef.co/go/tools/analysis/facts/generated"
	"honnef.co/go/tools/analysis/lint"
	"honnef.co/go/tools/unused"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"
	"github.com/golangci/golangci-lint/pkg/lint/linter"
	"github.com/golangci/golangci-lint/pkg/result"
)

const unusedName = "unused"

func NewUnused(settings *config.UnusedSettings, scSettings *config.StaticCheckSettings) *goanalysis.Linter {
	var mu sync.Mutex
	var resIssues []goanalysis.Issue

	analyzer := &analysis.Analyzer{
		Name:     unusedName,
		Doc:      unused.Analyzer.Analyzer.Doc,
		Requires: unused.Analyzer.Analyzer.Requires,
		Run: func(pass *analysis.Pass) (any, error) {
			issues := runUnused(pass, settings)
			if len(issues) == 0 {
				return nil, nil
			}

			mu.Lock()
			resIssues = append(resIssues, issues...)
			mu.Unlock()

			return nil, nil
		},
	}

	setAnalyzerGoVersion(analyzer, getGoVersion(scSettings))

	return goanalysis.NewLinter(
		unusedName,
		"Checks Go code for unused constants, variables, functions and types",
		[]*analysis.Analyzer{analyzer},
		nil,
	).WithIssuesReporter(func(_ *linter.Context) []goanalysis.Issue {
		return resIssues
	}).WithLoadMode(goanalysis.LoadModeTypesInfo)
}

func runUnused(pass *analysis.Pass, cfg *config.UnusedSettings) []goanalysis.Issue {
	res := getUnusedResults(pass, cfg)

	used := make(map[string]bool)
	for _, obj := range res.Used {
		used[fmt.Sprintf("%s %d %s", obj.Position.Filename, obj.Position.Line, obj.Name)] = true
	}

	var issues []goanalysis.Issue

	// Inspired by https://github.com/dominikh/go-tools/blob/d694aadcb1f50c2d8ac0a1dd06217ebb9f654764/lintcmd/lint.go#L177-L197
	for _, object := range res.Unused {
		if object.Kind == "type param" {
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

	return issues
}

func getUnusedResults(pass *analysis.Pass, settings *config.UnusedSettings) unused.Result {
	opts := unused.Options{
		FieldWritesAreUses:     settings.FieldWritesAreUses,
		PostStatementsAreReads: settings.PostStatementsAreReads,
		ExportedIsUsed:         settings.ExportedIsUsed,
		ExportedFieldsAreUsed:  settings.ExportedFieldsAreUsed,
		ParametersAreUsed:      settings.ParametersAreUsed,
		LocalVariablesAreUsed:  settings.LocalVariablesAreUsed,
		GeneratedIsUsed:        settings.GeneratedIsUsed,
	}

	// ref: https://github.com/dominikh/go-tools/blob/4ec1f474ca6c0feb8e10a8fcca4ab95f5b5b9881/internal/cmd/unused/unused.go#L68
	nodes := unused.Graph(pass.Fset,
		pass.Files,
		pass.Pkg,
		pass.TypesInfo,
		pass.ResultOf[directives.Analyzer].([]lint.Directive),
		pass.ResultOf[generated.Analyzer].(map[string]generated.Generator),
		opts,
	)

	sg := unused.SerializedGraph{}
	sg.Merge(nodes)
	return sg.Results()
}
