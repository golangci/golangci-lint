package goanalysis

import (
	"context"

	"github.com/pkg/errors"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/lint/linter"
	"github.com/golangci/golangci-lint/pkg/result"
)

type MetaLinter struct {
	linters              []*Linter
	analyzerToLinterName map[*analysis.Analyzer]string
}

func NewMetaLinter(linters []*Linter) *MetaLinter {
	ml := &MetaLinter{linters: linters}
	ml.analyzerToLinterName = ml.getAnalyzerToLinterNameMapping()
	return ml
}

func (ml MetaLinter) Name() string {
	return "goanalysis_metalinter"
}

func (ml MetaLinter) Desc() string {
	return ""
}

func (ml MetaLinter) isTypecheckMode() bool {
	for _, linter := range ml.linters {
		if linter.isTypecheckMode() {
			return true
		}
	}
	return false
}

func (ml MetaLinter) getLoadMode() LoadMode {
	loadMode := LoadModeNone
	for _, linter := range ml.linters {
		if linter.loadMode > loadMode {
			loadMode = linter.loadMode
		}
	}
	return loadMode
}

func (ml MetaLinter) getAnalyzers() []*analysis.Analyzer {
	var allAnalyzers []*analysis.Analyzer
	for _, linter := range ml.linters {
		allAnalyzers = append(allAnalyzers, linter.analyzers...)
	}
	return allAnalyzers
}

func (ml MetaLinter) getName() string {
	return "metalinter"
}

func (ml MetaLinter) useOriginalPackages() bool {
	return false // `unused` can't be run by this metalinter
}

func (ml MetaLinter) reportIssues(lintCtx *linter.Context) []Issue {
	var ret []Issue
	for _, lnt := range ml.linters {
		if lnt.issuesReporter != nil {
			ret = append(ret, lnt.issuesReporter(lintCtx)...)
		}
	}
	return ret
}

func (ml MetaLinter) getLinterNameForDiagnostic(diag *Diagnostic) string {
	return ml.analyzerToLinterName[diag.Analyzer]
}

func (ml MetaLinter) getAnalyzerToLinterNameMapping() map[*analysis.Analyzer]string {
	analyzerToLinterName := map[*analysis.Analyzer]string{}
	for _, linter := range ml.linters {
		for _, a := range linter.analyzers {
			analyzerToLinterName[a] = linter.Name()
		}
	}
	return analyzerToLinterName
}

func (ml MetaLinter) Run(ctx context.Context, lintCtx *linter.Context) ([]result.Issue, error) {
	for _, linter := range ml.linters {
		if err := linter.preRun(lintCtx); err != nil {
			return nil, errors.Wrapf(err, "failed to pre-run %s", linter.Name())
		}
	}

	return runAnalyzers(ml, lintCtx)
}
