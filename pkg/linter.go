package pkg

import (
	"context"

	"github.com/golangci/golangci-lint/pkg/golinters"
	"github.com/golangci/golangci-lint/pkg/result"
)

type Linter interface {
	Run(ctx context.Context, lintCtx *golinters.Context) ([]result.Issue, error)
	Name() string
}
