package golinters

import (
	"context"
	"fmt"
	"go/token"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"golang.org/x/tools/go/packages"

	"github.com/golangci/golangci-lint/pkg/lint/linter"
	libpackages "github.com/golangci/golangci-lint/pkg/packages"
	"github.com/golangci/golangci-lint/pkg/result"
)

type TypeCheck struct{}

func (TypeCheck) Name() string {
	return "typecheck"
}

func (TypeCheck) Desc() string {
	return "Like the front-end of a Go compiler, parses and type-checks Go code"
}

func (lint TypeCheck) parseError(srcErr packages.Error) (*result.Issue, error) {
	// file:line(<optional>:colon)
	parts := strings.Split(srcErr.Pos, ":")
	if len(parts) == 1 {
		return nil, errors.New("no colons")
	}

	file := parts[0]
	line, err := strconv.Atoi(parts[1])
	if err != nil {
		return nil, fmt.Errorf("can't parse line number %q: %s", parts[1], err)
	}

	var column int
	if len(parts) == 3 { // no column
		column, err = strconv.Atoi(parts[2])
		if err != nil {
			return nil, errors.Wrapf(err, "failed to parse column from %q", parts[2])
		}
	}

	return &result.Issue{
		Pos: token.Position{
			Filename: file,
			Line:     line,
			Column:   column,
		},
		Text:       srcErr.Msg,
		FromLinter: lint.Name(),
	}, nil
}

func (lint TypeCheck) Run(ctx context.Context, lintCtx *linter.Context) ([]result.Issue, error) {
	var res []result.Issue
	for _, pkg := range lintCtx.NotCompilingPackages {
		errors := libpackages.ExtractErrors(pkg)
		for _, err := range errors {
			i, perr := lint.parseError(err)
			if perr != nil { // failed to parse
				lintCtx.Log.Errorf("typechecking error: %s", err.Msg)
			} else {
				res = append(res, *i)
			}
		}
	}

	return res, nil
}
