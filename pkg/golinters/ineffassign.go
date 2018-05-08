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

func (Ineffassign) Desc() string {
	return "Detects when assignments to existing variables are not used"
}

func (lint Ineffassign) Run(ctx context.Context, lintCtx *Context) ([]result.Issue, error) {
	issues := ineffassignAPI.Run(lintCtx.Paths.Files)

	var res []result.Issue
	for _, i := range issues {
		res = append(res, result.Issue{
			Pos:        i.Pos,
			Text:       fmt.Sprintf("ineffectual assignment to %s", formatCode(i.IdentName, lintCtx.RunCfg())),
			FromLinter: lint.Name(),
		})
	}
	return res, nil
}
