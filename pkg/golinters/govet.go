package golinters

import (
	"context"

	"github.com/golangci/golangci-lint/pkg/result"
	govetAPI "github.com/golangci/govet"
)

type Govet struct{}

func (Govet) Name() string {
	return "govet"
}

func (g Govet) Run(ctx context.Context, lintCtx *Context) ([]result.Issue, error) {
	issues, err := govetAPI.Run(lintCtx.Paths.MixedPaths(), lintCtx.RunCfg().BuildTags, lintCtx.RunCfg().Govet.CheckShadowing)
	if err != nil {
		return nil, err
	}

	var res []result.Issue
	for _, i := range issues {
		res = append(res, result.Issue{
			Pos:        i.Pos,
			Text:       i.Message,
			FromLinter: g.Name(),
		})
	}
	return res, nil
}
