package golinters

import (
	"context"
	"fmt"

	"github.com/golangci/golangci-lint/pkg/result"
	ineffassignAPI "github.com/golangci/ineffassign"
)

type Ineffassign struct{}

func (Ineffassign) Name() string {
	return "ineffassign"
}

func (lint Ineffassign) Run(ctx context.Context, lintCtx *Context) (*result.Result, error) {
	issues := ineffassignAPI.Run(lintCtx.Paths.Files)

	res := &result.Result{}
	for _, i := range issues {
		res.Issues = append(res.Issues, result.Issue{
			File:       i.Pos.Filename,
			LineNumber: i.Pos.Line,
			Text:       fmt.Sprintf("ineffectual assignment to %s", formatCode(i.IdentName, lintCtx.RunCfg())),
			FromLinter: lint.Name(),
		})
	}
	return res, nil
}
