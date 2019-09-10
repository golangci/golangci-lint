package golinters

import (
	"context"
	"go/token"

	"github.com/golangci/golangci-lint/pkg/lint/linter"
	"github.com/golangci/golangci-lint/pkg/result"

	"github.com/ultraware/whitespace"
)

type Whitespace struct{}

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
		res[k] = result.Issue{
			Pos: token.Position{
				Filename: i.Pos.Filename,
				Line:     i.Pos.Line,
			},
			Text:       i.Message,
			FromLinter: w.Name(),
		}
	}

	return res, nil
}
