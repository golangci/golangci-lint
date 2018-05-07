package golinters

import (
	"context"

	"mvdan.cc/interfacer/check"

	"github.com/golangci/golangci-lint/pkg/result"
)

type Interfacer struct{}

func (Interfacer) Name() string {
	return "interfacer"
}

func (lint Interfacer) Run(ctx context.Context, lintCtx *Context) (*result.Result, error) {
	c := new(check.Checker)
	c.Program(lintCtx.Program)
	c.ProgramSSA(lintCtx.SSAProgram)

	issues, err := c.Check()
	if err != nil {
		return nil, err
	}

	res := &result.Result{}
	for _, i := range issues {
		pos := lintCtx.SSAProgram.Fset.Position(i.Pos())
		res.Issues = append(res.Issues, result.Issue{
			File:       pos.Filename,
			LineNumber: pos.Line,
			Text:       i.Message(),
			FromLinter: lint.Name(),
		})
	}

	return res, nil
}
