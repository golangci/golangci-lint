package golinters

import (
	"context"
	"fmt"

	goconstAPI "github.com/golangci/goconst"
	"github.com/golangci/golangci-lint/pkg/result"
)

type Goconst struct{}

func (Goconst) Name() string {
	return "goconst"
}

func (lint Goconst) Run(ctx context.Context, lintCtx *Context) (*result.Result, error) {
	issues, err := goconstAPI.Run(lintCtx.Paths.Files, true,
		lintCtx.RunCfg().Goconst.MinStringLen,
		lintCtx.RunCfg().Goconst.MinOccurrencesCount,
	)
	if err != nil {
		return nil, err
	}

	res := &result.Result{}
	for _, i := range issues {
		textBegin := fmt.Sprintf("string %s has %d occurrences", formatCode(i.Str, lintCtx.RunCfg()), i.OccurencesCount)
		var textEnd string
		if i.MatchingConst == "" {
			textEnd = ", make it a constant"
		} else {
			textEnd = fmt.Sprintf(", but such constant %s already exists", formatCode(i.MatchingConst, lintCtx.RunCfg()))
		}
		res.Issues = append(res.Issues, result.Issue{
			File:       i.Pos.Filename,
			LineNumber: i.Pos.Line,
			Text:       textBegin + textEnd,
			FromLinter: lint.Name(),
		})
	}

	return res, nil
}
