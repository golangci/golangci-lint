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
	settings := whitespace.Settings{MultiIf: lintCtx.Cfg.LintersSettings.Whitespace.MultiIf}

	var issues []whitespace.Message
	for _, file := range lintCtx.ASTCache.GetAllValidFiles() {
		issues = append(issues, whitespace.Run(file.F, file.Fset, settings)...)
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
			LineRange:   &result.Range{From: i.Pos.Line, To: i.Pos.Line},
			Text:        i.Message,
			FromLinter:  w.Name(),
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

		res[k] = issue
	}

	return res, nil
}
