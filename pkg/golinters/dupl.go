package golinters

import (
	"context"
	"fmt"

	"github.com/golangci/golangci-lint/pkg/result"
	duplAPI "github.com/mibk/dupl"
)

type Dupl struct{}

func (Dupl) Name() string {
	return "dupl"
}

func (d Dupl) Run(ctx context.Context, lintCtx *Context) ([]result.Issue, error) {
	issues, err := duplAPI.Run(lintCtx.Paths.Files, lintCtx.RunCfg().Dupl.Threshold)
	if err != nil {
		return nil, err
	}

	var res []result.Issue
	for _, i := range issues {
		dupl := fmt.Sprintf("%s:%d-%d", i.To.Filename(), i.To.LineStart(), i.To.LineEnd())
		text := fmt.Sprintf("%d-%d lines are duplicate of %s",
			i.From.LineStart(), i.From.LineEnd(),
			formatCode(dupl, lintCtx.RunCfg()))
		res = append(res, result.Issue{
			File:       i.From.Filename(),
			LineNumber: i.From.LineStart(),
			Text:       text,
			FromLinter: d.Name(),
		})
	}
	return res, nil
}
