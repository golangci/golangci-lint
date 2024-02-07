package golinters

import (
	"fmt"
	"sync"

	"github.com/ultraware/whitespace"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"
	"github.com/golangci/golangci-lint/pkg/lint/linter"
	"github.com/golangci/golangci-lint/pkg/result"
)

const whitespaceName = "whitespace"

func NewWhitespace(settings *config.WhitespaceSettings) *goanalysis.Linter {
	var mu sync.Mutex
	var resIssues []goanalysis.Issue

	var wsSettings whitespace.Settings
	if settings != nil {
		wsSettings = whitespace.Settings{
			Mode:      whitespace.RunningModeGolangCI,
			MultiIf:   settings.MultiIf,
			MultiFunc: settings.MultiFunc,
		}
	}

	a := whitespace.NewAnalyzer(&wsSettings)

	return goanalysis.NewLinter(
		a.Name,
		a.Doc,
		[]*analysis.Analyzer{a},
		nil,
	).WithContextSetter(func(_ *linter.Context) {
		a.Run = func(pass *analysis.Pass) (any, error) {
			issues, err := runWhitespace(pass, wsSettings)
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

func runWhitespace(pass *analysis.Pass, wsSettings whitespace.Settings) ([]goanalysis.Issue, error) {
	lintIssues := whitespace.Run(pass, &wsSettings)

	issues := make([]goanalysis.Issue, len(lintIssues))
	for i, issue := range lintIssues {
		report := &result.Issue{
			FromLinter: whitespaceName,
			Pos:        pass.Fset.PositionFor(issue.Diagnostic, false),
			Text:       issue.Message,
		}

		switch issue.MessageType {
		case whitespace.MessageTypeRemove:
			if len(issue.LineNumbers) == 0 {
				continue
			}

			report.LineRange = &result.Range{
				From: issue.LineNumbers[0],
				To:   issue.LineNumbers[len(issue.LineNumbers)-1],
			}

			report.Replacement = &result.Replacement{NeedOnlyDelete: true}

		case whitespace.MessageTypeAdd:
			report.Pos = pass.Fset.PositionFor(issue.FixStart, false)
			report.Replacement = &result.Replacement{
				Inline: &result.InlineFix{
					StartCol:  0,
					Length:    1,
					NewString: "\n\t",
				},
			}

		default:
			return nil, fmt.Errorf("unknown message type: %v", issue.MessageType)
		}

		issues[i] = goanalysis.NewIssue(report, pass)
	}

	return issues, nil
}
