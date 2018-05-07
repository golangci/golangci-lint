package golinters

import (
	"context"

	"github.com/golangci/golangci-lint/pkg/result"
	unconvertAPI "github.com/golangci/unconvert"
)

type Unconvert struct{}

func (Unconvert) Name() string {
	return "unconvert"
}

func (lint Unconvert) Run(ctx context.Context, lintCtx *Context) (*result.Result, error) {
	positions := unconvertAPI.Run(lintCtx.Program)
	res := &result.Result{}
	for _, pos := range positions {
		res.Issues = append(res.Issues, result.Issue{
			File:       pos.Filename,
			LineNumber: pos.Line,
			Text:       "unnecessary conversion",
			FromLinter: lint.Name(),
		})
	}

	return res, nil
}
