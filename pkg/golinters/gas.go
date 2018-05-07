package golinters

import (
	"context"
	"fmt"
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

func (lint Gas) Run(ctx context.Context, lintCtx *Context) (*result.Result, error) {
	gasConfig := gas.NewConfig()
	enabledRules := rules.Generate(rules.NewRuleFilter(true, "G104")) // disable what errcheck does: it reports on Close etc
	logger := log.New(ioutil.Discard, "", 0)
	analyzer := gas.NewAnalyzer(gasConfig, logger)
	analyzer.LoadRules(enabledRules.Builders())

	analyzer.ProcessProgram(lintCtx.Program)
	issues, _ := analyzer.Report()

	res := &result.Result{}
	for _, i := range issues {
		text := fmt.Sprintf("%s: %s", i.RuleID, i.What) // TODO: use severity and confidence
		line, _ := strconv.Atoi(i.Line)
		res.Issues = append(res.Issues, result.Issue{
			File:       i.File,
			LineNumber: line,
			Text:       text,
			FromLinter: lint.Name(),
		})
	}

	return res, nil
}
