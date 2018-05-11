package golinters // nolint:dupl

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

func (Structcheck) Desc() string {
	return "Finds unused struct fields"
}

func (s Structcheck) Run(ctx context.Context, lintCtx *Context) ([]result.Issue, error) {
	issues := structcheckAPI.Run(lintCtx.Program, lintCtx.Settings().Structcheck.CheckExportedFields)

	var res []result.Issue
	for _, i := range issues {
		res = append(res, result.Issue{
			Pos:        i.Pos,
			Text:       fmt.Sprintf("%s is unused", formatCode(i.FieldName, lintCtx.Cfg)),
			FromLinter: s.Name(),
		})
	}
	return res, nil
}
