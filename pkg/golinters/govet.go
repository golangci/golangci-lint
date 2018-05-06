package golinters

import (
	"context"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/result"
	"github.com/golangci/golangci-shared/pkg/executors"
	govetAPI "github.com/golangci/govet"
)

type govet struct{}

func (govet) Name() string {
	return "govet"
}

func (g govet) Run(ctx context.Context, exec executors.Executor, cfg *config.Run) (*result.Result, error) {
	issues, err := govetAPI.Run(cfg.Paths.MixedPaths(), cfg.BuildTags, cfg.Govet.CheckShadowing)
	if err != nil {
		return nil, err
	}

	res := &result.Result{}
	for _, i := range issues {
		res.Issues = append(res.Issues, result.Issue{
			File:       i.Pos.Filename,
			LineNumber: i.Pos.Line,
			Text:       i.Message,
			FromLinter: g.Name(),
		})
	}
	return res, nil
}
