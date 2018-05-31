package golinters

import (
	"context"
	"fmt"

	"github.com/golangci/golangci-lint/pkg/lint"
	"github.com/golangci/golangci-lint/pkg/result"
	errcheckAPI "github.com/kisielk/errcheck/golangci"
)

type Errcheck struct{}

func (Errcheck) Name() string {
	return "errcheck"
}

func (Errcheck) Desc() string {
	return "Errcheck is a program for checking for unchecked errors in go programs. These unchecked errors can be critical bugs in some cases"
}

func (e Errcheck) Run(ctx context.Context, lintCtx *lint.Context) ([]result.Issue, error) {
	errCfg := &lintCtx.Settings().Errcheck
	issues, err := errcheckAPI.Run(lintCtx.Program, errCfg.CheckAssignToBlank, errCfg.CheckTypeAssertions)
	if err != nil {
		return nil, err
	}

	if len(issues) == 0 {
		return nil, nil
	}

	res := make([]result.Issue, 0, len(issues))
	for _, i := range issues {
		var text string
		if i.FuncName != "" {
			text = fmt.Sprintf("Error return value of %s is not checked", formatCode(i.FuncName, lintCtx.Cfg))
		} else {
			text = "Error return value is not checked"
		}
		res = append(res, result.Issue{
			FromLinter: e.Name(),
			Text:       text,
			Pos:        i.Pos,
		})
	}

	return res, nil
}
