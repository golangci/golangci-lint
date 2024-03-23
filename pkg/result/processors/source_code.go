package processors

import (
	"github.com/golangci/golangci-lint/pkg/fsutils"
	"github.com/golangci/golangci-lint/pkg/logutils"
	"github.com/golangci/golangci-lint/pkg/result"
)

var _ Processor = (*SourceCode)(nil)

type SourceCode struct {
	lineCache *fsutils.LineCache
	log       logutils.Log
}

func NewSourceCode(lc *fsutils.LineCache, log logutils.Log) *SourceCode {
	return &SourceCode{
		lineCache: lc,
		log:       log,
	}
}

func (SourceCode) Name() string {
	return "source_code"
}

func (p SourceCode) Process(issues []result.Issue) ([]result.Issue, error) {
	return transformIssues(issues, p.transform), nil
}

func (SourceCode) Finish() {}

func (p SourceCode) transform(issue *result.Issue) *result.Issue {
	newIssue := *issue

	lineRange := issue.GetLineRange()
	for lineNumber := lineRange.From; lineNumber <= lineRange.To; lineNumber++ {
		line, err := p.lineCache.GetLine(issue.FilePath(), lineNumber)
		if err != nil {
			p.log.Warnf("Failed to get line %d for file %s: %s",
				lineNumber, issue.FilePath(), err)

			return issue
		}

		newIssue.SourceLines = append(newIssue.SourceLines, line)
	}

	return &newIssue
}
