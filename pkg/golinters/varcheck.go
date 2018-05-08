package golinters

import (
	"context"
	"fmt"

	varcheckAPI "github.com/golangci/check/cmd/varcheck"
	"github.com/golangci/golangci-lint/pkg/result"
)

type Varcheck struct{}

func (Varcheck) Name() string {
	return "varcheck"
}

func (v Varcheck) Run(ctx context.Context, lintCtx *Context) ([]result.Issue, error) {
	issues := varcheckAPI.Run(lintCtx.Program, lintCtx.RunCfg().Varcheck.CheckExportedFields)

	var res []result.Issue
	for _, i := range issues {
		res = append(res, result.Issue{
			Pos:        i.Pos,
			Text:       fmt.Sprintf("%s is unused", formatCode(i.VarName, lintCtx.RunCfg())),
			FromLinter: v.Name(),
		})
	}
	return res, nil
}
