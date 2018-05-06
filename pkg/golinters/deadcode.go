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

func (d Deadcode) Run(ctx context.Context, lintCtx *Context) (*result.Result, error) {
	issues, err := deadcodeAPI.Run(lintCtx.Program)
	if err != nil {
		return nil, err
	}

	res := &result.Result{}
	for _, i := range issues {
		res.Issues = append(res.Issues, result.Issue{
			File:       i.Pos.Filename,
			LineNumber: i.Pos.Line,
			Text:       fmt.Sprintf("%s is unused", formatCode(i.UnusedIdentName, lintCtx.RunCfg())),
			FromLinter: d.Name(),
		})
	}
	return res, nil
}
