package golinters

import (
	"context"
	"go/token"
	"strconv"
	"strings"

	"github.com/golangci/golangci-lint/pkg/result"
)

type TypeCheck struct{}

func (TypeCheck) Name() string {
	return "typecheck"
}

func (TypeCheck) Desc() string {
	return "Like the front-end of a Go compiler, parses and type-checks Go code"
}

func (lint TypeCheck) parseError(err error) *result.Issue {
	// file:line(<optional>:colon): message
	parts := strings.Split(err.Error(), ":")
	if len(parts) < 3 {
		return nil
	}

	file := parts[0]
	line, err := strconv.Atoi(parts[1])
	if err != nil {
		return nil
	}

	var column int
	var message string
	if len(parts) == 3 { // no column
		message = parts[2]
	} else {
		column, err = strconv.Atoi(parts[2])
		if err == nil { // column was parsed
			message = strings.Join(parts[3:], ":")
		} else {
			message = strings.Join(parts[2:], ":")
		}
	}

	message = strings.TrimSpace(message)
	if message == "" {
		return nil
	}

	return &result.Issue{
		Pos: token.Position{
			Filename: file,
			Line:     line,
			Column:   column,
		},
		Text:       message,
		FromLinter: lint.Name(),
	}
}

func (lint TypeCheck) Run(ctx context.Context, lintCtx *Context) ([]result.Issue, error) {
	if lintCtx.NotCompilingPackages == nil {
		return nil, nil
	}

	var res []result.Issue
	for _, pkg := range lintCtx.NotCompilingPackages {
		for _, err := range pkg.Errors {
			i := lint.parseError(err)
			if i != nil {
				res = append(res, *i)
			}
		}
	}

	return res, nil
}
