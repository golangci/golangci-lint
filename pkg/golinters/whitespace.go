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

//nolint:dupl
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

	whitespaceAnalyzer := whitespace.NewAnalyzer(&wsSettings)

	return goanalysis.NewLinter(
		whitespaceAnalyzer.Name,
		whitespaceAnalyzer.Doc,
		[]*analysis.Analyzer{whitespaceAnalyzer},
		nil,
	).WithContextSetter(func(lintCtx *linter.Context) {
		whitespaceAnalyzer.Run = func(pass *analysis.Pass) (any, error) {
			whitespaceIssues := whitespace.Run(pass, &wsSettings)
			issues := make([]goanalysis.Issue, len(whitespaceIssues))

			for i, issue := range whitespaceIssues {
				report := &result.Issue{
					FromLinter: whitespaceAnalyzer.Name,
					Pos:        pass.Fset.PositionFor(issue.Diagnostic, false),
					Text:       issue.Message,
				}

				switch issue.MessageType {
				case whitespace.MessageTypeRemove:
					report.LineRange = &result.Range{
						From: issue.LineNumber,
						To:   issue.LineNumber,
					}

					report.Replacement = &result.Replacement{
						NeedOnlyDelete: true,
					}

				case whitespace.MessageTypeAdd:
					position := pass.Fset.PositionFor(issue.Diagnostic, false)
					bracketLine, err := lintCtx.LineCache.GetLine(position.Filename, position.Line)
					if err != nil {
						return nil, fmt.Errorf("failed to get line %s:%d: %w", position.Filename, position.Line, err)
					}

					report.Replacement = &result.Replacement{
						NewLines: []string{bracketLine + "\n"},
					}

				default:
					return nil, fmt.Errorf("unknown message type: %v", issue.MessageType)
				}

				issues[i] = goanalysis.NewIssue(report, pass)
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
