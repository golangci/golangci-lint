package golinters

import (
	"go/token"
	"sync"

	"github.com/pkg/errors"
	"github.com/ultraware/whitespace"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"
	"github.com/golangci/golangci-lint/pkg/lint/linter"
	"github.com/golangci/golangci-lint/pkg/result"
)

const whitespaceName = "whitespace"

//nolint:dupl
func NewWhitespace(settings *config.WhitespaceSettings) *goanalysis.Linter {
	var mu sync.Mutex
	var resIssues []goanalysis.Issue

	var wsSettings whitespace.Settings
	if settings != nil {
		wsSettings = whitespace.Settings{
			MultiIf:   settings.MultiIf,
			MultiFunc: settings.MultiFunc,
		}
	}

	analyzer := &analysis.Analyzer{
		Name: whitespaceName,
		Doc:  goanalysis.TheOnlyanalyzerDoc,
		Run:  goanalysis.DummyRun,
	}

	return goanalysis.NewLinter(
		whitespaceName,
		"Tool for detection of leading and trailing whitespace",
		[]*analysis.Analyzer{analyzer},
		nil,
	).WithContextSetter(func(lintCtx *linter.Context) {
		analyzer.Run = func(pass *analysis.Pass) (interface{}, error) {
			issues, err := runWhitespace(lintCtx, pass, wsSettings)
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

func runWhitespace(lintCtx *linter.Context, pass *analysis.Pass, wsSettings whitespace.Settings) ([]goanalysis.Issue, error) {
	var messages []whitespace.Message
	for _, file := range pass.Files {
		messages = append(messages, whitespace.Run(file, pass.Fset, wsSettings)...)
	}

	if len(messages) == 0 {
		return nil, nil
	}

	issues := make([]goanalysis.Issue, len(messages))
	for k, i := range messages {
		issue := result.Issue{
			Pos: token.Position{
				Filename: i.Pos.Filename,
				Line:     i.Pos.Line,
			},
			LineRange:   &result.Range{From: i.Pos.Line, To: i.Pos.Line},
			Text:        i.Message,
			FromLinter:  whitespaceName,
			Replacement: &result.Replacement{},
		}

		bracketLine, err := lintCtx.LineCache.GetLine(issue.Pos.Filename, issue.Pos.Line)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to get line %s:%d", issue.Pos.Filename, issue.Pos.Line)
		}

		switch i.Type {
		case whitespace.MessageTypeLeading:
			issue.LineRange.To++ // cover two lines by the issue: opening bracket "{" (issue.Pos.Line) and following empty line
		case whitespace.MessageTypeTrailing:
			issue.LineRange.From-- // cover two lines by the issue: closing bracket "}" (issue.Pos.Line) and preceding empty line
			issue.Pos.Line--       // set in sync with LineRange.From to not break fixer and other code features
		case whitespace.MessageTypeAddAfter:
			bracketLine += "\n"
		}
		issue.Replacement.NewLines = []string{bracketLine}

		issues[k] = goanalysis.NewIssue(&issue, pass)
	}

	return issues, nil
}
