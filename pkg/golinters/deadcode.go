package golinters

import (
	"context"
	"fmt"

	deadcodeAPI "github.com/golangci/go-misc/deadcode"
	"github.com/golangci/golangci-lint/pkg/result"
)

type Deadcode struct{}

func (Deadcode) Name() string {
	return "deadcode"
}

func (d Deadcode) Run(ctx context.Context, lintCtx *Context) ([]result.Issue, error) {
	issues, err := deadcodeAPI.Run(lintCtx.Program)
	if err != nil {
		return nil, err
	}

	var res []result.Issue
	for _, i := range issues {
		res = append(res, result.Issue{
			Pos:        i.Pos,
			Text:       fmt.Sprintf("%s is unused", formatCode(i.UnusedIdentName, lintCtx.RunCfg())),
			FromLinter: d.Name(),
		})
	}
	return res, nil
}
