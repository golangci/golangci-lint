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
	"github.com/sirupsen/logrus"
)

type SimpleRunner struct {
	Processors []processors.Processor
}

type lintRes struct {
	linter Linter
	err    error
	issues []result.Issue
}

func (r *SimpleRunner) runLinter(ctx context.Context, linter Linter, lintCtx *golinters.Context, i int) (res []result.Issue, err error) {
	defer func() {
		if panicData := recover(); panicData != nil {
			err = fmt.Errorf("panic occured: %s", panicData)
			logrus.Infof("Panic stack trace: %s", debug.Stack())
		}
	}()
	startedAt := time.Now()
	res, err = linter.Run(ctx, lintCtx)

	logrus.Infof("worker #%d: linter %s took %s and found %d issues (before processing them)", i, linter.Name(),
		time.Since(startedAt), len(res))
	return
}

func (r *SimpleRunner) runLinters(ctx context.Context, wg *sync.WaitGroup, tasksCh chan Linter, lintResultsCh chan lintRes, lintCtx *golinters.Context, workersCount int) {
	for i := 0; i < workersCount; i++ {
		go func(i int) {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case linter, ok := <-tasksCh:
					if !ok {
						return
					}
					if ctx.Err() != nil {
						// XXX: if check it in only int a select
						// it's possible to not enter to this case until tasksCh is empty.
						return
					}
					issues, lerr := r.runLinter(ctx, linter, lintCtx, i)
					lintResultsCh <- lintRes{
						linter: linter,
						err:    lerr,
						issues: issues,
					}
				}
			}
		}(i + 1)
	}
}

func (r SimpleRunner) Run(ctx context.Context, linters []Linter, lintCtx *golinters.Context) chan result.Issue {
	retIssues := make(chan result.Issue, 1024)
	go func() {
		defer close(retIssues)
		if err := r.runGo(ctx, linters, lintCtx, retIssues); err != nil {
			logrus.Warnf("error running linters: %s", err)
		}
	}()

	return retIssues
}

func (r SimpleRunner) runGo(ctx context.Context, linters []Linter, lintCtx *golinters.Context, retIssues chan result.Issue) error {
	savedStdout, savedStderr := os.Stdout, os.Stderr
	devNull, err := os.Open(os.DevNull)
	if err != nil {
		return fmt.Errorf("can't open null device %q: %s", os.DevNull, err)
	}

	// Don't allow linters to print anything
	os.Stdout, os.Stderr = devNull, devNull

	lintResultsCh := make(chan lintRes, len(linters))
	tasksCh := make(chan Linter, len(linters))
	workersCount := lintCtx.Cfg.Run.Concurrency
	var wg sync.WaitGroup
	wg.Add(workersCount)
	r.runLinters(ctx, &wg, tasksCh, lintResultsCh, lintCtx, workersCount)

	for _, linter := range linters {
		tasksCh <- linter
	}
	close(tasksCh)

	go func() {
		wg.Wait()
		close(lintResultsCh)
		os.Stdout, os.Stderr = savedStdout, savedStderr
	}()

	finishedN := 0
	for res := range lintResultsCh {
		if res.err != nil {
			logrus.Warnf("Can't run linter %s: %s", res.linter.Name(), res.err)
			continue
		}

		finishedN++

		if len(res.issues) != 0 {
			res.issues = r.processIssues(ctx, res.issues)
		}

		for _, i := range res.issues {
			retIssues <- i
		}
	}

	// finalize processors: logging, clearing, no heavy work here
	for _, p := range r.Processors {
		p.Finish()
	}

	if ctx.Err() != nil {
		return fmt.Errorf("%d/%d linters finished: deadline exceeded: try increase it by passing --deadline option",
			finishedN, len(linters))
	}

	return nil
}

func (r *SimpleRunner) processIssues(ctx context.Context, issues []result.Issue) []result.Issue {
	for _, p := range r.Processors {
		startedAt := time.Now()
		newIssues, err := p.Process(issues)
		elapsed := time.Since(startedAt)
		if elapsed > 50*time.Millisecond {
			logrus.Infof("Result processor %s took %s", p.Name(), elapsed)
		}
		if err != nil {
			logrus.Warnf("Can't process result by %s processor: %s", p.Name(), err)
		} else {
			issues = newIssues
		}
		if issues == nil {
			issues = []result.Issue{}
		}
	}

	return issues
}
