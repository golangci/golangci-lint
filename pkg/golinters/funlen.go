package golinters

import (
	"context"
	"go/token"

	"github.com/golangci/golangci-lint/pkg/lint/linter"
	"github.com/golangci/golangci-lint/pkg/result"

	"github.com/ultraware/funlen"
)

type Funlen struct{}

func (Funlen) Name() string {
	return "funlen"
}

func (Funlen) Desc() string {
	return "Tool for detection of long functions"
}

func (f Funlen) Run(ctx context.Context, lintCtx *linter.Context) ([]result.Issue, error) {
	var issues []funlen.Message
	for _, file := range lintCtx.ASTCache.GetAllValidFiles() {
		issues = append(issues, funlen.Run(file.F, file.Fset, lintCtx.Settings().Funlen.Lines, lintCtx.Settings().Funlen.Statements)...)
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
			FromLinter: f.Name(),
		}
	}

	return res, nil
}
