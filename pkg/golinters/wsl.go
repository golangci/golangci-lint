package golinters

import (
	"sync"

	"github.com/bombsimon/wsl/v3"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"
	"github.com/golangci/golangci-lint/pkg/lint/linter"
	"github.com/golangci/golangci-lint/pkg/result"
)

const wslName = "wsl"

// NewWSL returns a new WSL linter.
func NewWSL(settings *config.WSLSettings) *goanalysis.Linter {
	var mu sync.Mutex
	var resIssues []goanalysis.Issue

	conf := &wsl.Configuration{
		AllowCuddleWithCalls: []string{"Lock", "RLock"},
		AllowCuddleWithRHS:   []string{"Unlock", "RUnlock"},
		ErrorVariableNames:   []string{"err"},
	}

	if settings != nil {
		conf.StrictAppend = settings.StrictAppend
		conf.AllowAssignAndCallCuddle = settings.AllowAssignAndCallCuddle
		conf.AllowAssignAndAnythingCuddle = settings.AllowAssignAndAnythingCuddle
		conf.AllowMultiLineAssignCuddle = settings.AllowMultiLineAssignCuddle
		conf.AllowCuddleDeclaration = settings.AllowCuddleDeclaration
		conf.AllowTrailingComment = settings.AllowTrailingComment
		conf.AllowSeparatedLeadingComment = settings.AllowSeparatedLeadingComment
		conf.ForceCuddleErrCheckAndAssign = settings.ForceCuddleErrCheckAndAssign
		conf.ForceCaseTrailingWhitespaceLimit = settings.ForceCaseTrailingWhitespaceLimit
		conf.ForceExclusiveShortDeclarations = settings.ForceExclusiveShortDeclarations
	}

	analyzer := &analysis.Analyzer{
		Name: goanalysis.TheOnlyAnalyzerName,
		Doc:  goanalysis.TheOnlyanalyzerDoc,
		Run: func(pass *analysis.Pass) (interface{}, error) {
			issues := runWSL(pass, conf)

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
		wslName,
		"Whitespace Linter - Forces you to use empty lines!",
		[]*analysis.Analyzer{analyzer},
		nil,
	).WithIssuesReporter(func(*linter.Context) []goanalysis.Issue {
		return resIssues
	}).WithLoadMode(goanalysis.LoadModeSyntax)
}

func runWSL(pass *analysis.Pass, conf *wsl.Configuration) []goanalysis.Issue {
	var files = make([]string, 0, len(pass.Files))
	for _, file := range pass.Files {
		files = append(files, pass.Fset.PositionFor(file.Pos(), false).Filename)
	}

	if conf == nil {
		return nil
	}

	wslErrors, _ := wsl.NewProcessorWithConfig(*conf).ProcessFiles(files)
	if len(wslErrors) == 0 {
		return nil
	}

	var issues []goanalysis.Issue
	for _, err := range wslErrors {
		issues = append(issues, goanalysis.NewIssue(&result.Issue{
			FromLinter: wslName,
			Pos:        err.Position,
			Text:       err.Reason,
		}, pass))
	}

	return issues
}
