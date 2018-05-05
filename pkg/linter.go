package linters

import (
	"context"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/result"
	"github.com/golangci/golangci-shared/pkg/executors"
)

type Linter interface {
	Run(ctx context.Context, exec executors.Executor, cfg *config.Run) (*result.Result, error)
	Name() string
}
