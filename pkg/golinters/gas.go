package golinters

import (
	"context"
	"fmt"
	"go/token"
	"io/ioutil"
	"log"
	"strconv"

	"github.com/GoASTScanner/gas"
	"github.com/GoASTScanner/gas/rules"
	"github.com/golangci/golangci-lint/pkg/result"
)

type Gas struct{}

func (Gas) Name() string {
	return "gas"
}

func (lint Gas) Run(ctx context.Context, lintCtx *Context) ([]result.Issue, error) {
	gasConfig := gas.NewConfig()
	enabledRules := rules.Generate()
	logger := log.New(ioutil.Discard, "", 0)
	analyzer := gas.NewAnalyzer(gasConfig, logger)
	analyzer.LoadRules(enabledRules.Builders())

	analyzer.ProcessProgram(lintCtx.Program)
	issues, _ := analyzer.Report()

	var res []result.Issue
	for _, i := range issues {
		text := fmt.Sprintf("%s: %s", i.RuleID, i.What) // TODO: use severity and confidence
		line, _ := strconv.Atoi(i.Line)
		res = append(res, result.Issue{
			Pos: token.Position{
				Filename: i.File,
				Line:     line,
			},
			Text:       text,
			FromLinter: lint.Name(),
		})
	}

	return res, nil
}
