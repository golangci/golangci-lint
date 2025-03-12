package goanalysis

import (
	"context"
	"fmt"

	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/v2/pkg/lint/linter"
	"github.com/golangci/golangci-lint/v2/pkg/result"
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

func (ml MetaLinter) Run(_ context.Context, lintCtx *linter.Context) ([]result.Issue, error) {
	for _, l := range ml.linters {
		if err := l.preRun(lintCtx); err != nil {
			return nil, fmt.Errorf("failed to pre-run %s: %w", l.Name(), err)
		}
	}

	return runAnalyzers(ml, lintCtx)
}

func (MetaLinter) Name() string {
	return "goanalysis_metalinter"
}

func (MetaLinter) Desc() string {
	return ""
}

func (ml MetaLinter) getLoadMode() LoadMode {
	loadMode := LoadModeNone
	for _, l := range ml.linters {
		if l.loadMode > loadMode {
			loadMode = l.loadMode
		}
	}
	return loadMode
}

func (ml MetaLinter) getAnalyzers() []*analysis.Analyzer {
	var allAnalyzers []*analysis.Analyzer
	for _, l := range ml.linters {
		allAnalyzers = append(allAnalyzers, l.analyzers...)
	}
	return allAnalyzers
}

func (MetaLinter) getName() string {
	return "metalinter"
}

func (MetaLinter) useOriginalPackages() bool {
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
	for _, l := range ml.linters {
		for _, a := range l.analyzers {
			analyzerToLinterName[a] = l.Name()
		}
	}
	return analyzerToLinterName
}
