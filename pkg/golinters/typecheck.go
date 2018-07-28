package golinters

import (
	"context"
	"errors"
	"fmt"
	"go/token"
	"strconv"
	"strings"

	"github.com/golangci/golangci-lint/pkg/lint/linter"
	"github.com/golangci/golangci-lint/pkg/result"
)

type TypeCheck struct{}

func (TypeCheck) Name() string {
	return "typecheck"
}

func (TypeCheck) Desc() string {
	return "Like the front-end of a Go compiler, parses and type-checks Go code"
}

func (lint TypeCheck) parseError(srcErr error) (*result.Issue, error) {
	// TODO: cast srcErr to types.Error and just use it

	// file:line(<optional>:colon): message
	parts := strings.Split(srcErr.Error(), ":")
	if len(parts) < 3 {
		return nil, errors.New("too few colons")
	}

	file := parts[0]
	line, err := strconv.Atoi(parts[1])
	if err != nil {
		return nil, fmt.Errorf("can't parse line number %q: %s", parts[1], err)
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
		return nil, fmt.Errorf("empty message")
	}

	return &result.Issue{
		Pos: token.Position{
			Filename: file,
			Line:     line,
			Column:   column,
		},
		Text:       message,
		FromLinter: lint.Name(),
	}, nil
}

func (lint TypeCheck) Run(ctx context.Context, lintCtx *linter.Context) ([]result.Issue, error) {
	var res []result.Issue
	for _, pkg := range lintCtx.NotCompilingPackages {
		for _, err := range pkg.Errors {
			i, perr := lint.parseError(err)
			if perr != nil {
				lintCtx.Log.Warnf("Can't parse type error %s: %s", err, perr)
			} else {
				res = append(res, *i)
			}
		}
	}

	return res, nil
}
