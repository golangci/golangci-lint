package pkg

import (
	"context"
	"fmt"
	"os"
	"runtime/debug"
	"sync"
	"time"

	"github.com/golangci/golangci-lint/pkg/golinters"
	"github.com/golangci/golangci-lint/pkg/result"
	"github.com/golangci/golangci-lint/pkg/result/processors"
	"github.com/golangci/golangci-shared/pkg/analytics"
)

type Runner interface {
	Run(ctx context.Context, linters []Linter, lintCtx *golinters.Context) ([]result.Issue, error)
}

type SimpleRunner struct {
	Processors []processors.Processor
}

type lintRes struct {
	linter Linter
	err    error
	res    *result.Result
}

func runLinter(ctx context.Context, linter Linter, lintCtx *golinters.Context, i int) (res *result.Result, err error) {
	defer func() {
		if panicData := recover(); panicData != nil {
			err = fmt.Errorf("panic occured: %s", panicData)
			analytics.Log(ctx).Infof("Panic stack trace: %s", debug.Stack())
		}
	}()
	startedAt := time.Now()
	res, err = linter.Run(ctx, lintCtx)
	analytics.Log(ctx).Infof("worker #%d: linter %s took %s for paths %s", i, linter.Name(),
		time.Since(startedAt), lintCtx.Paths.MixedPaths())
	return
}

func runLinters(ctx context.Context, wg *sync.WaitGroup, tasksCh chan Linter, lintResultsCh chan lintRes, lintCtx *golinters.Context) {
	for i := 0; i < lintCtx.Cfg.Common.Concurrency; i++ {
		go func(i int) {
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
					res, lerr := runLinter(ctx, linter, lintCtx, i)
					lintResultsCh <- lintRes{
						linter: linter,
						err:    lerr,
						res:    res,
					}
				}
			}
		}(i + 1)
	}
}

func (r SimpleRunner) Run(ctx context.Context, linters []Linter, lintCtx *golinters.Context) ([]result.Issue, error) {
	savedStdout, savedStderr := os.Stdout, os.Stderr
	devNull, err := os.Open(os.DevNull)
	if err != nil {
		return nil, fmt.Errorf("can't open null device %q: %s", os.DevNull, err)
	}

	os.Stdout, os.Stderr = devNull, devNull

	lintResultsCh := make(chan lintRes, len(linters))
	tasksCh := make(chan Linter, lintCtx.Cfg.Common.Concurrency)
	var wg sync.WaitGroup
	wg.Add(lintCtx.Cfg.Common.Concurrency)
	runLinters(ctx, &wg, tasksCh, lintResultsCh, lintCtx)

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

	results, err = r.processResults(ctx, results)
	if err != nil {
		return nil, fmt.Errorf("can't process results: %s", err)
	}

	return r.mergeResults(results), nil
}

func (r SimpleRunner) processResults(ctx context.Context, results []result.Result) ([]result.Result, error) {
	if len(r.Processors) == 0 {
		return results, nil
	}

	for _, p := range r.Processors {
		var err error
		startedAt := time.Now()
		results, err = p.Process(results)
		elapsed := time.Since(startedAt)
		if elapsed > 50*time.Millisecond {
			analytics.Log(ctx).Infof("Result processor %s took %s", p.Name(), elapsed)
		}
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
