package golinters

import (
	"context"
	"fmt"

	goconstAPI "github.com/golangci/goconst"
	"github.com/golangci/golangci-lint/pkg/lint/linter"
	"github.com/golangci/golangci-lint/pkg/result"
)

type Goconst struct{}

func (Goconst) Name() string {
	return "goconst"
}

func (Goconst) Desc() string {
	return "Finds repeated strings that could be replaced by a constant"
}

func (lint Goconst) Run(ctx context.Context, lintCtx *linter.Context) ([]result.Issue, error) {
	var goconstIssues []goconstAPI.Issue
	// TODO: make it cross-package: pass package names inside goconst
	for _, files := range lintCtx.Paths.FilesGrouppedByDirs() {
		issues, err := goconstAPI.Run(files, true,
			lintCtx.Settings().Goconst.MinStringLen,
			lintCtx.Settings().Goconst.MinOccurrencesCount,
		)
		if err != nil {
			return nil, err
		}

		goconstIssues = append(goconstIssues, issues...)
	}
	if len(goconstIssues) == 0 {
		return nil, nil
	}

	res := make([]result.Issue, 0, len(goconstIssues))
	for _, i := range goconstIssues {
		textBegin := fmt.Sprintf("string %s has %d occurrences", formatCode(i.Str, lintCtx.Cfg), i.OccurencesCount)
		var textEnd string
		if i.MatchingConst == "" {
			textEnd = ", make it a constant"
		} else {
			textEnd = fmt.Sprintf(", but such constant %s already exists", formatCode(i.MatchingConst, lintCtx.Cfg))
		}
		res = append(res, result.Issue{
			Pos:        i.Pos,
			Text:       textBegin + textEnd,
			FromLinter: lint.Name(),
		})
	}

	return res, nil
}
