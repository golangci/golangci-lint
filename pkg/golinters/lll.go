package golinters

import (
	"bufio"
	"context"
	"fmt"
	"go/token"
	"os"
	"strings"
	"unicode/utf8"

	"github.com/golangci/golangci-lint/pkg/lint/linter"
	"github.com/golangci/golangci-lint/pkg/result"
)

type Lll struct{}

func (Lll) Name() string {
	return "lll"
}

func (Lll) Desc() string {
	return "Reports long lines"
}

func (lint Lll) getIssuesForFile(filename string, maxLineLen int, tabSpaces string) ([]result.Issue, error) {
	var res []result.Issue

	f, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("can't open file %s: %s", filename, err)
	}
	defer f.Close()

	lineNumber := 1
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.Replace(line, "\t", tabSpaces, -1)
		lineLen := utf8.RuneCountInString(line)
		if lineLen > maxLineLen {
			res = append(res, result.Issue{
				Pos: token.Position{
					Filename: filename,
					Line:     lineNumber,
					Column:   1,
				},
				Text:       fmt.Sprintf("line is %d characters", lineLen),
				FromLinter: lint.Name(),
			})
		}
		lineNumber++
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("can't scan file %s: %s", filename, err)
	}

	return res, nil
}

func (lint Lll) Run(ctx context.Context, lintCtx *linter.Context) ([]result.Issue, error) {
	var res []result.Issue
	spaces := strings.Repeat(" ", lintCtx.Settings().Lll.TabWidth)
	for _, f := range lintCtx.PkgProgram.Files(lintCtx.Cfg.Run.AnalyzeTests) {
		issues, err := lint.getIssuesForFile(f, lintCtx.Settings().Lll.LineLength, spaces)
		if err != nil {
			return nil, err
		}
		res = append(res, issues...)
	}

	return res, nil
}
