package goanalysis

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/lint/linter"
	"github.com/golangci/golangci-lint/pkg/result"
)

type MetaLinter struct {
	linters              []*Linter
	analyzerToLinterName map[*analysis.Analyzer]string
}

func NewMetaLinter(linters []*Linter, analyzerToLinterName map[*analysis.Analyzer]string) *MetaLinter {
	return &MetaLinter{linters: linters, analyzerToLinterName: analyzerToLinterName}
}

func (ml MetaLinter) Name() string {
	return "goanalysis_metalinter"
}

func (ml MetaLinter) Desc() string {
	return ""
}

func (ml MetaLinter) Run(ctx context.Context, lintCtx *linter.Context) ([]result.Issue, error) {
	for _, linter := range ml.linters {
		if err := analysis.Validate(linter.analyzers); err != nil {
			return nil, errors.Wrapf(err, "failed to validate analyzers of %s", linter.Name())
		}
	}

	for _, linter := range ml.linters {
		if err := linter.configure(); err != nil {
			return nil, errors.Wrapf(err, "failed to configure analyzers of %s", linter.Name())
		}
	}

	var allAnalyzers []*analysis.Analyzer
	for _, linter := range ml.linters {
		allAnalyzers = append(allAnalyzers, linter.analyzers...)
	}

	runner := newRunner("metalinter", lintCtx.Log.Child("goanalysis"), lintCtx.PkgCache, lintCtx.LoadGuard, lintCtx.NeedWholeProgram)

	diags, errs := runner.run(allAnalyzers, lintCtx.Packages)
	// Don't print all errs: they can duplicate.
	if len(errs) != 0 {
		return nil, errs[0]
	}

	var issues []result.Issue
	for i := range diags {
		diag := &diags[i]
		issues = append(issues, result.Issue{
			FromLinter: ml.analyzerToLinterName[diag.Analyzer],
			Text:       fmt.Sprintf("%s: %s", diag.Analyzer, diag.Message),
			Pos:        diag.Position,
		})
	}

	return issues, nil
}
