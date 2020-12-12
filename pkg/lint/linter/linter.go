package linter

import (
	"context"

	"github.com/anduril/golangci-lint/pkg/result"
)

type Linter interface {
	Run(ctx context.Context, lintCtx *Context) ([]result.Issue, error)
	Name() string
	Desc() string
}
