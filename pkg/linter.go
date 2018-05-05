package linters

import (
	"context"

	"github.com/golangci/golangci-lint/pkg/result"
	"github.com/golangci/golangci-shared/pkg/executors"
)

type Linter interface {
	Run(ctx context.Context, exec executors.Executor) (*result.Result, error)
	Name() string
}
