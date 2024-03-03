package linter

import (
	"context"

	"github.com/golangci/golangci-lint/pkg/config"
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
}

func NewNoop(l Linter, reason string) Noop {
	return Noop{
		name:   l.Name(),
		desc:   l.Desc(),
		reason: reason,
	}
}

func NewNoopDeprecated(name string, cfg *config.Config) Noop {
	noop := Noop{
		name:   name,
		desc:   "Deprecated",
		reason: "This linter is fully inactivated: it will not produce any reports.",
	}

	if cfg.InternalCmdTest {
		noop.reason = ""
	}

	return noop
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
