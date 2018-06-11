package golinters

import (
	"context"
	"fmt"
	"go/token"

	"github.com/golangci/golangci-lint/pkg/lint/linter"
	"github.com/golangci/golangci-lint/pkg/result"
	duplAPI "github.com/mibk/dupl"
)

type Dupl struct{}

func (Dupl) Name() string {
	return "dupl"
}

func (Dupl) Desc() string {
	return "Tool for code clone detection"
}

func (d Dupl) Run(ctx context.Context, lintCtx *linter.Context) ([]result.Issue, error) {
	issues, err := duplAPI.Run(lintCtx.PkgProgram.Files(lintCtx.Cfg.Run.AnalyzeTests), lintCtx.Settings().Dupl.Threshold)
	if err != nil {
		return nil, err
	}

	if len(issues) == 0 {
		return nil, nil
	}

	res := make([]result.Issue, 0, len(issues))
	for _, i := range issues {
		dupl := fmt.Sprintf("%s:%d-%d", i.To.Filename(), i.To.LineStart(), i.To.LineEnd())
		text := fmt.Sprintf("%d-%d lines are duplicate of %s",
			i.From.LineStart(), i.From.LineEnd(),
			formatCode(dupl, lintCtx.Cfg))
		res = append(res, result.Issue{
			Pos: token.Position{
				Filename: i.From.Filename(),
				Line:     i.From.LineStart(),
			},
			LineRange: &result.Range{
				From: i.From.LineStart(),
				To:   i.From.LineEnd(),
			},
			Text:       text,
			FromLinter: d.Name(),
		})
	}
	return res, nil
}
