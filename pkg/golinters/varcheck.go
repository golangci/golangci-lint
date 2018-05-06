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

func (v Varcheck) Run(ctx context.Context, lintCtx *Context) (*result.Result, error) {
	issues := varcheckAPI.Run(lintCtx.Program, lintCtx.RunCfg().Varcheck.CheckExportedFields)

	res := &result.Result{}
	for _, i := range issues {
		res.Issues = append(res.Issues, result.Issue{
			File:       i.Pos.Filename,
			LineNumber: i.Pos.Line,
			Text:       fmt.Sprintf("%s is unused", formatCode(i.VarName, lintCtx.RunCfg())),
			FromLinter: v.Name(),
		})
	}
	return res, nil
}
