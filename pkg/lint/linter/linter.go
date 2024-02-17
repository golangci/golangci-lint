package linter

import (
	"context"

	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/result"
)

type Linter interface {
	Run(ctx context.Context, lintCtx *Context) ([]result.Issue, error)
	Name() string
	Desc() string
}

type Noop struct {
	name   string
	desc   string
	reason string
	run    func(pass *analysis.Pass) (any, error)
}

func (n Noop) Run(_ context.Context, lintCtx *Context) ([]result.Issue, error) {
	if n.reason != "" {
		lintCtx.Log.Warnf("%s: %s", n.name, n.reason)
	}
	return nil, nil
}

func (n Noop) Name() string {
	return n.name
}

func (n Noop) Desc() string {
	return n.desc
}
