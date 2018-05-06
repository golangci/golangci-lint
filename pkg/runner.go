package pkg

import (
	"context"
	"fmt"
	"os"
	"sync"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/result"
	"github.com/golangci/golangci-lint/pkg/result/processors"
	"github.com/golangci/golangci-shared/pkg/analytics"
	"github.com/golangci/golangci-shared/pkg/executors"
)

type Runner interface {
	Run(ctx context.Context, linters []Linter, exec executors.Executor, cfg *config.Config) ([]result.Issue, error)
}

type SimpleRunner struct {
	Processors []processors.Processor
}

type lintRes struct {
	linter Linter
	err    error
	res    *result.Result
}

func runLinter(ctx context.Context, linter Linter, exec executors.Executor, cfg *config.Run) (res *result.Result, err error) {
	defer func() {
		if panicData := recover(); panicData != nil {
			err = fmt.Errorf("panic occured: %s", panicData)
		}
	}()
	res, err = linter.Run(ctx, exec, cfg)
	return
}

func runLinters(ctx context.Context, wg *sync.WaitGroup, tasksCh chan Linter, lintResultsCh chan lintRes, exec executors.Executor, cfg *config.Config) {
	for i := 0; i < cfg.Common.Concurrency; i++ {
		go func() {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					// XXX: if check it in a select with reading from tasksCh
					// it's possible to not enter to this case until tasksCh is empty.
					return
				default:
				}

				select {
				case <-ctx.Done():
					return
				case linter, ok := <-tasksCh:
					if !ok {
						return
					}
					res, lerr := runLinter(ctx, linter, exec, &cfg.Run)
					lintResultsCh <- lintRes{
						linter: linter,
						err:    lerr,
						res:    res,
					}
				}
			}
		}()
	}
}

func (r SimpleRunner) Run(ctx context.Context, linters []Linter, exec executors.Executor, cfg *config.Config) ([]result.Issue, error) {
	savedStdout, savedStderr := os.Stdout, os.Stderr
	devNull, err := os.Open(os.DevNull)
	if err != nil {
		return nil, fmt.Errorf("can't open null device %q: %s", os.DevNull, err)
	}

	os.Stdout, os.Stderr = devNull, devNull

	lintResultsCh := make(chan lintRes, len(linters))
	tasksCh := make(chan Linter, cfg.Common.Concurrency)
	var wg sync.WaitGroup
	wg.Add(cfg.Common.Concurrency)
	runLinters(ctx, &wg, tasksCh, lintResultsCh, exec, cfg)

	for _, linter := range linters {
		tasksCh <- linter
	}

	close(tasksCh)
	wg.Wait()
	close(lintResultsCh)

	os.Stdout, os.Stderr = savedStdout, savedStderr
	results := []result.Result{}
	for res := range lintResultsCh {
		if res.err != nil {
			analytics.Log(ctx).Warnf("Can't run linter %s: %s", res.linter.Name(), res.err)
			continue
		}

		if res.res == nil || len(res.res.Issues) == 0 {
			continue
		}

		results = append(results, *res.res)
	}

	results, err = r.processResults(results)
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
