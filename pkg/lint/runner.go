package lint

import (
	"context"
	"fmt"
	"runtime/debug"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/golangci/golangci-lint/pkg/lint/linter"
	"github.com/golangci/golangci-lint/pkg/logutils"
	"github.com/golangci/golangci-lint/pkg/result"
	"github.com/golangci/golangci-lint/pkg/result/processors"
	"github.com/golangci/golangci-lint/pkg/timeutils"
	"github.com/sirupsen/logrus"
)

type SimpleRunner struct {
	Processors []processors.Processor
}

type lintRes struct {
	linter linter.Config
	err    error
	issues []result.Issue
}

func runLinterSafe(ctx context.Context, lintCtx *linter.Context, lc linter.Config) (ret []result.Issue, err error) {
	defer func() {
		if panicData := recover(); panicData != nil {
			err = fmt.Errorf("panic occured: %s", panicData)
			logutils.HiddenWarnf("Panic stack trace: %s", debug.Stack())
		}
	}()

	issues, err := lc.Linter.Run(ctx, lintCtx)
	if err != nil {
		return nil, err
	}

	for _, i := range issues {
		i.FromLinter = lc.Linter.Name()
	}

	return issues, nil
}

func runWorker(ctx context.Context, lintCtx *linter.Context, tasksCh <-chan linter.Config, lintResultsCh chan<- lintRes, name string) {
	sw := timeutils.NewStopwatch(name)
	defer sw.Print()

	for {
		select {
		case <-ctx.Done():
			return
		case lc, ok := <-tasksCh:
			if !ok {
				return
			}
			if ctx.Err() != nil {
				// XXX: if check it in only int a select
				// it's possible to not enter to this case until tasksCh is empty.
				return
			}
			var issues []result.Issue
			var err error
			sw.TrackStage(lc.Linter.Name(), func() {
				issues, err = runLinterSafe(ctx, lintCtx, lc)
			})
			lintResultsCh <- lintRes{
				linter: lc,
				err:    err,
				issues: issues,
			}
		}
	}
}

func logWorkersStat(workersFinishTimes []time.Time) {
	lastFinishTime := workersFinishTimes[0]
	for _, t := range workersFinishTimes {
		if t.After(lastFinishTime) {
			lastFinishTime = t
		}
	}

	logStrings := []string{}
	for i, t := range workersFinishTimes {
		if t.Equal(lastFinishTime) {
			continue
		}

		logStrings = append(logStrings, fmt.Sprintf("#%d: %s", i+1, lastFinishTime.Sub(t)))
	}

	logrus.Infof("Workers idle times: %s", strings.Join(logStrings, ", "))
}

func getSortedLintersConfigs(linters []linter.Config) []linter.Config {
	ret := make([]linter.Config, len(linters))
	copy(ret, linters)

	sort.Slice(ret, func(i, j int) bool {
		return ret[i].GetSpeed() < ret[j].GetSpeed()
	})

	return ret
}

func (r *SimpleRunner) runWorkers(ctx context.Context, lintCtx *linter.Context, linters []linter.Config) <-chan lintRes {
	tasksCh := make(chan linter.Config, len(linters))
	lintResultsCh := make(chan lintRes, len(linters))
	var wg sync.WaitGroup

	workersFinishTimes := make([]time.Time, lintCtx.Cfg.Run.Concurrency)

	for i := 0; i < lintCtx.Cfg.Run.Concurrency; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			name := fmt.Sprintf("worker.%d", i+1)
			runWorker(ctx, lintCtx, tasksCh, lintResultsCh, name)
			workersFinishTimes[i] = time.Now()
		}(i)
	}

	lcs := getSortedLintersConfigs(linters)
	for _, lc := range lcs {
		tasksCh <- lc
	}
	close(tasksCh)

	go func() {
		wg.Wait()
		close(lintResultsCh)

		logWorkersStat(workersFinishTimes)
	}()

	return lintResultsCh
}

func (r SimpleRunner) processLintResults(ctx context.Context, inCh <-chan lintRes) <-chan lintRes {
	outCh := make(chan lintRes, 64)

	go func() {
		sw := timeutils.NewStopwatch("processing")

		defer close(outCh)

		for res := range inCh {
			if res.err != nil {
				logutils.HiddenWarnf("Can't run linter %s: %s", res.linter.Linter.Name(), res.err)
				continue
			}

			if len(res.issues) != 0 {
				res.issues = r.processIssues(ctx, res.issues, sw)
				outCh <- res
			}
		}

		// finalize processors: logging, clearing, no heavy work here

		for _, p := range r.Processors {
			sw.TrackStage(p.Name(), func() {
				p.Finish()
			})
		}

		sw.PrintStages()
	}()

	return outCh
}

func collectIssues(ctx context.Context, resCh <-chan lintRes) <-chan result.Issue {
	retIssues := make(chan result.Issue, 1024)
	go func() {
		defer close(retIssues)

		for res := range resCh {
			if len(res.issues) == 0 {
				continue
			}

			for _, i := range res.issues {
				retIssues <- i
			}
		}
	}()

	return retIssues
}

func (r SimpleRunner) Run(ctx context.Context, linters []linter.Config, lintCtx *linter.Context) <-chan result.Issue {
	lintResultsCh := r.runWorkers(ctx, lintCtx, linters)
	processedLintResultsCh := r.processLintResults(ctx, lintResultsCh)
	if ctx.Err() != nil {
		// XXX: always process issues, even if timeout occured
		finishedLintersN := 0
		for range processedLintResultsCh {
			finishedLintersN++
		}

		logrus.Warnf("%d/%d linters finished: deadline exceeded: try increase it by passing --deadline option",
			finishedLintersN, len(linters))
	}

	return collectIssues(ctx, processedLintResultsCh)
}

func (r *SimpleRunner) processIssues(ctx context.Context, issues []result.Issue, sw *timeutils.Stopwatch) []result.Issue {
	for _, p := range r.Processors {
		var newIssues []result.Issue
		var err error
		sw.TrackStage(p.Name(), func() {
			newIssues, err = p.Process(issues)
		})

		if err != nil {
			logutils.HiddenWarnf("Can't process result by %s processor: %s", p.Name(), err)
		} else {
			issues = newIssues
		}

		if issues == nil {
			issues = []result.Issue{}
		}
	}

	return issues
}
