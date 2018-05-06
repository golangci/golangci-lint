package golinters

import (
	"context"
	"fmt"

	deadcodeAPI "github.com/golangci/go-misc/deadcode"
	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/result"
	"github.com/golangci/golangci-shared/pkg/executors"
)

type deadcode struct{}

func (deadcode) Name() string {
	return "deadcode"
}

func (d deadcode) Run(ctx context.Context, exec executors.Executor, cfg *config.Run) (*result.Result, error) {
	issues, err := deadcodeAPI.Run(cfg.Paths.MixedPaths(), true) // TODO: configure need of tests
	if err != nil {
		return nil, err
	}

	res := &result.Result{}
	for _, i := range issues {
		res.Issues = append(res.Issues, result.Issue{
			File:       i.Pos.Filename,
			LineNumber: i.Pos.Line,
			Text:       fmt.Sprintf("%s is unused", formatCode(i.UnusedIdentName, cfg)),
			FromLinter: d.Name(),
		})
	}
	return res, nil
}
