package golinters

import (
	"context"
	"fmt"

	structcheckAPI "github.com/golangci/check/cmd/structcheck"
	"github.com/golangci/golangci-lint/pkg/result"
)

type Structcheck struct{}

func (Structcheck) Name() string {
	return "structcheck"
}

func (s Structcheck) Run(ctx context.Context, lintCtx *Context) (*result.Result, error) {
	issues := structcheckAPI.Run(lintCtx.Program, lintCtx.RunCfg().Structcheck.CheckExportedFields)

	res := &result.Result{}
	for _, i := range issues {
		res.Issues = append(res.Issues, result.Issue{
			File:       i.Pos.Filename,
			LineNumber: i.Pos.Line,
			Text:       fmt.Sprintf("%s is unused", formatCode(i.FieldName, lintCtx.RunCfg())),
			FromLinter: s.Name(),
		})
	}
	return res, nil
}
