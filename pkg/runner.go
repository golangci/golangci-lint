package pkg

import (
	"context"
	"fmt"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/result"
	"github.com/golangci/golangci-lint/pkg/result/processors"
	"github.com/golangci/golangci-shared/pkg/analytics"
	"github.com/golangci/golangci-shared/pkg/executors"
)

type Runner interface {
	Run(ctx context.Context, linters []Linter, exec executors.Executor, cfg *config.Run) ([]result.Issue, error)
}

type SimpleRunner struct {
	Processors []processors.Processor
}

func (r SimpleRunner) Run(ctx context.Context, linters []Linter, exec executors.Executor, cfg *config.Run) ([]result.Issue, error) {
	results := []result.Result{}
	for _, linter := range linters {
		res, err := linter.Run(ctx, exec, cfg)
		if err != nil {
			analytics.Log(ctx).Warnf("Can't run linter %+v: %s", linter, err)
			continue
		}

		if len(res.Issues) == 0 {
			continue
		}

		results = append(results, *res)
	}

	results, err := r.processResults(results)
	if err != nil {
		return nil, fmt.Errorf("can't process results: %s", err)
	}

	return r.mergeResults(results), nil
}

func (r SimpleRunner) processResults(results []result.Result) ([]result.Result, error) {
	if len(r.Processors) == 0 {
		return results, nil
	}

	for _, p := range r.Processors {
		var err error
		results, err = p.Process(results)
		if err != nil {
			return nil, err
		}
	}

	return results, nil
}

func (r SimpleRunner) mergeResults(results []result.Result) []result.Issue {
	issues := []result.Issue{}
	for _, r := range results {
		issues = append(issues, r.Issues...)
	}

	return issues
}
