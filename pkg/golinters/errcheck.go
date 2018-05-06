package golinters

import (
	"context"
	"fmt"
	"strings"

	"github.com/golangci/golangci-lint/pkg/result"
	errcheckAPI "github.com/kisielk/errcheck/golangci"
)

type Errcheck struct{}

func (Errcheck) Name() string {
	return "errcheck"
}

func (e Errcheck) Run(ctx context.Context, lintCtx *Context) (*result.Result, error) {
	errCfg := &lintCtx.RunCfg().Errcheck
	issues, err := errcheckAPI.Run(lintCtx.Program, errCfg.CheckAssignToBlank, errCfg.CheckTypeAssertions)
	if err != nil {
		return nil, err
	}

	res := &result.Result{}
	for _, i := range issues {
		if !errCfg.CheckClose && strings.HasSuffix(i.FuncName, ".Close") {
			continue
		}

		var text string
		if i.FuncName != "" {
			text = fmt.Sprintf("Error return value of %s is not checked", formatCode(i.FuncName, lintCtx.RunCfg()))
		} else {
			text = "Error return value is not checked"
		}
		res.Issues = append(res.Issues, result.Issue{
			FromLinter: e.Name(),
			Text:       text,
			LineNumber: i.Pos.Line,
			File:       i.Pos.Filename,
		})
	}

	return res, nil
}
