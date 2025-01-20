package lint

import (
	"context"
	"errors"
	"fmt"
	"runtime/debug"
	"strings"

	"github.com/golangci/golangci-lint/internal/errorutil"
	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/fsutils"
	"github.com/golangci/golangci-lint/pkg/goformatters"
	"github.com/golangci/golangci-lint/pkg/goutil"
	"github.com/golangci/golangci-lint/pkg/lint/linter"
	"github.com/golangci/golangci-lint/pkg/lint/lintersdb"
	"github.com/golangci/golangci-lint/pkg/logutils"
	"github.com/golangci/golangci-lint/pkg/result"
	"github.com/golangci/golangci-lint/pkg/result/processors"
	"github.com/golangci/golangci-lint/pkg/timeutils"
)

type processorStat struct {
	inCount  int
	outCount int
}

type Runner struct {
	Log logutils.Log

	lintCtx    *linter.Context
	Processors []processors.Processor
}

func NewRunner(log logutils.Log, cfg *config.Config, args []string, goenv *goutil.Env,
	lineCache *fsutils.LineCache, fileCache *fsutils.FileCache,
	dbManager *lintersdb.Manager, lintCtx *linter.Context,
) (*Runner, error) {
	// Beware that some processors need to add the path prefix when working with paths
	// because they get invoked before the path prefixer (exclude and severity rules)
	// or process other paths (skip files).
	files := fsutils.NewFiles(lineCache, cfg.Output.PathPrefix)

	pathRelativity, err := processors.NewPathRelativity(log, cfg.GetBasePath())
	if err != nil {
		return nil, fmt.Errorf("error creating path relativity processor: %w", err)
	}

	exclusionPaths, err := processors.NewExclusionPaths(log, &cfg.Linters.LinterExclusions)
	if err != nil {
		return nil, err
	}

	skipFilesProcessor, err := processors.NewSkipFiles(cfg.Issues.ExcludeFiles, cfg.Output.PathPrefix)
	if err != nil {
		return nil, err
	}

	skipDirs := cfg.Issues.ExcludeDirs
	if cfg.Issues.UseDefaultExcludeDirs {
		skipDirs = append(skipDirs, processors.StdExcludeDirRegexps...)
	}

	skipDirsProcessor, err := processors.NewSkipDirs(log.Child(logutils.DebugKeySkipDirs), skipDirs, args, cfg.Output.PathPrefix)
	if err != nil {
		return nil, err
	}

	enabledLinters, err := dbManager.GetEnabledLintersMap()
	if err != nil {
		return nil, fmt.Errorf("failed to get enabled linters: %w", err)
	}

	metaFormatter, err := goformatters.NewMetaFormatter(log, cfg, enabledLinters)
	if err != nil {
		return nil, fmt.Errorf("failed to create meta-formatter: %w", err)
	}

	return &Runner{
		Processors: []processors.Processor{
			// Must be the first processor.
			processors.NewPathAbsoluter(log),

			processors.NewCgo(goenv),

			// Must be after Cgo.
			processors.NewFilenameUnadjuster(lintCtx.Packages, log.Child(logutils.DebugKeyFilenameUnadjuster)),

			// Must be after FilenameUnadjuster.
			processors.NewInvalidIssue(log.Child(logutils.DebugKeyInvalidIssue)),

			// Must be after PathAbsoluter, Cgo, FilenameUnadjuster InvalidIssue.
			pathRelativity,

			// Must be after PathRelativity.
			exclusionPaths,
			skipFilesProcessor,
			skipDirsProcessor,

			processors.NewGeneratedFileFilter(cfg.Linters.LinterExclusions.Generated),

			// Must be before exclude because users see already marked output and configure excluding by it.
			processors.NewIdentifierMarker(),

			processors.NewExclusionRules(log.Child(logutils.DebugKeyExclusionRules), files,
				&cfg.Linters.LinterExclusions, &cfg.Issues),

			processors.NewNolintFilter(log.Child(logutils.DebugKeyNolintFilter), dbManager, enabledLinters),

			processors.NewDiff(&cfg.Issues),

			// The fixer still needs to see paths for the issues that are relative to the current directory.
			processors.NewFixer(cfg, log, fileCache, metaFormatter),

			// Must be after the Fixer.
			processors.NewUniqByLine(cfg.Issues.UniqByLine),
			processors.NewMaxPerFileFromLinter(cfg),
			processors.NewMaxSameIssues(cfg.Issues.MaxSameIssues, log.Child(logutils.DebugKeyMaxSameIssues), cfg),
			processors.NewMaxFromLinter(cfg.Issues.MaxIssuesPerLinter, log.Child(logutils.DebugKeyMaxFromLinter), cfg),

			// Now we can modify the issues for output.
			processors.NewSourceCode(lineCache, log.Child(logutils.DebugKeySourceCode)),
			processors.NewPathShortener(),
			processors.NewSeverity(log.Child(logutils.DebugKeySeverityRules), files, &cfg.Severity),
			processors.NewPathPrettifier(log, cfg.Output.PathPrefix),
			processors.NewSortResults(&cfg.Output),
		},
		lintCtx: lintCtx,
		Log:     log,
	}, nil
}

func (r *Runner) Run(ctx context.Context, linters []*linter.Config) ([]result.Issue, error) {
	sw := timeutils.NewStopwatch("linters", r.Log)
	defer sw.Print()

	var (
		lintErrors error
		issues     []result.Issue
	)

	for _, lc := range linters {
		linterIssues, err := timeutils.TrackStage(sw, lc.Name(), func() ([]result.Issue, error) {
			return r.runLinterSafe(ctx, r.lintCtx, lc)
		})
		if err != nil {
			lintErrors = errors.Join(lintErrors, fmt.Errorf("can't run linter %s", lc.Linter.Name()), err)
			r.Log.Warnf("Can't run linter %s: %v", lc.Linter.Name(), err)

			continue
		}

		issues = append(issues, linterIssues...)
	}

	return r.processLintResults(issues), lintErrors
}

func (r *Runner) runLinterSafe(ctx context.Context, lintCtx *linter.Context,
	lc *linter.Config,
) (ret []result.Issue, err error) {
	defer func() {
		if panicData := recover(); panicData != nil {
			if pe, ok := panicData.(*errorutil.PanicError); ok {
				err = fmt.Errorf("%s: %w", lc.Name(), pe)

				// Don't print stacktrace from goroutines twice
				r.Log.Errorf("Panic: %s: %s", pe, pe.Stack())
			} else {
				err = fmt.Errorf("panic occurred: %s", panicData)
				r.Log.Errorf("Panic stack trace: %s", debug.Stack())
			}
		}
	}()

	issues, err := lc.Linter.Run(ctx, lintCtx)

	if lc.DoesChangeTypes {
		// Packages in lintCtx might be dirty due to the last analysis,
		// which affects to the next analysis.
		// To avoid this issue, we clear type information from the packages.
		// See https://github.com/golangci/golangci-lint/pull/944.
		// Currently, DoesChangeTypes is true only for `unused`.
		lintCtx.ClearTypesInPackages()
	}

	if err != nil {
		return nil, err
	}

	for i := range issues {
		if issues[i].FromLinter == "" {
			issues[i].FromLinter = lc.Name()
		}
	}

	return issues, nil
}

func (r *Runner) processLintResults(inIssues []result.Issue) []result.Issue {
	sw := timeutils.NewStopwatch("processing", r.Log)

	var issuesBefore, issuesAfter int
	statPerProcessor := map[string]processorStat{}

	var outIssues []result.Issue
	if len(inIssues) != 0 {
		issuesBefore += len(inIssues)
		outIssues = r.processIssues(inIssues, sw, statPerProcessor)
		issuesAfter += len(outIssues)
	}

	// finalize processors: logging, clearing, no heavy work here

	for _, p := range r.Processors {
		sw.TrackStage(p.Name(), p.Finish)
	}

	if issuesBefore != issuesAfter {
		r.Log.Infof("Issues before processing: %d, after processing: %d", issuesBefore, issuesAfter)
	}
	r.printPerProcessorStat(statPerProcessor)
	sw.PrintStages()

	return outIssues
}

func (r *Runner) printPerProcessorStat(stat map[string]processorStat) {
	parts := make([]string, 0, len(stat))
	for name, ps := range stat {
		if ps.inCount != 0 {
			parts = append(parts, fmt.Sprintf("%s: %d/%d", name, ps.inCount, ps.outCount))
		}
	}
	if len(parts) != 0 {
		r.Log.Infof("Processors filtering stat (in/out): %s", strings.Join(parts, ", "))
	}
}

func (r *Runner) processIssues(issues []result.Issue, sw *timeutils.Stopwatch, statPerProcessor map[string]processorStat) []result.Issue {
	for _, p := range r.Processors {
		newIssues, err := timeutils.TrackStage(sw, p.Name(), func() ([]result.Issue, error) {
			return p.Process(issues)
		})

		if err != nil {
			r.Log.Warnf("Can't process result by %s processor: %s", p.Name(), err)
		} else {
			stat := statPerProcessor[p.Name()]
			stat.inCount += len(issues)
			stat.outCount += len(newIssues)
			statPerProcessor[p.Name()] = stat
			issues = newIssues
		}

		// This is required by JSON serialization
		if issues == nil {
			issues = []result.Issue{}
		}
	}

	return issues
}
