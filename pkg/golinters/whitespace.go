package golinters

import (
	"context"
	"go/token"

	"github.com/pkg/errors"

	"github.com/golangci/golangci-lint/pkg/lint/linter"
	"github.com/golangci/golangci-lint/pkg/result"

	"github.com/ultraware/whitespace"
)

type Whitespace struct {
}

func (Whitespace) Name() string {
	return "whitespace"
}

func (Whitespace) Desc() string {
	return "Tool for detection of leading and trailing whitespace"
}

func (w Whitespace) Run(ctx context.Context, lintCtx *linter.Context) ([]result.Issue, error) {
	var issues []whitespace.Message
	for _, file := range lintCtx.ASTCache.GetAllValidFiles() {
		issues = append(issues, whitespace.Run(file.F, file.Fset)...)
	}

	if len(issues) == 0 {
		return nil, nil
	}

	res := make([]result.Issue, len(issues))
	for k, i := range issues {
		issue := result.Issue{
			Pos: token.Position{
				Filename: i.Pos.Filename,
				Line:     i.Pos.Line,
			},
			Text:        i.Message,
			FromLinter:  w.Name(),
			Replacement: &result.Replacement{},
		}

		// TODO(jirfag): return more information from Whitespace to get rid of string comparisons
		if i.Message == "unnecessary leading newline" {
			// cover two lines by the issue: opening bracket "{" (issue.Pos.Line) and following empty line
			issue.LineRange = &result.Range{From: issue.Pos.Line, To: issue.Pos.Line + 1}

			bracketLine, err := lintCtx.LineCache.GetLine(issue.Pos.Filename, issue.Pos.Line)
			if err != nil {
				return nil, errors.Wrapf(err, "failed to get line %s:%d", issue.Pos.Filename, issue.Pos.Line)
			}
			issue.Replacement.NewLines = []string{bracketLine}
		} else {
			// cover two lines by the issue: closing bracket "}" (issue.Pos.Line) and preceding empty line
			issue.LineRange = &result.Range{From: issue.Pos.Line - 1, To: issue.Pos.Line}

			bracketLine, err := lintCtx.LineCache.GetLine(issue.Pos.Filename, issue.Pos.Line)
			if err != nil {
				return nil, errors.Wrapf(err, "failed to get line %s:%d", issue.Pos.Filename, issue.Pos.Line)
			}
			issue.Replacement.NewLines = []string{bracketLine}

			issue.Pos.Line-- // set in sync with LineRange.From to not break fixer and other code features
		}

		res[k] = issue
	}

	return res, nil
}
