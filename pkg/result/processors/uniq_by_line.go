package processors

import (
	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/result"
)

const uniqByLineLimit = 1

var _ Processor = (*UniqByLine)(nil)

type UniqByLine struct {
	fileLineCounter fileLineCounter
	cfg             *config.Config
}

func NewUniqByLine(cfg *config.Config) *UniqByLine {
	return &UniqByLine{
		fileLineCounter: fileLineCounter{},
		cfg:             cfg,
	}
}

func (*UniqByLine) Name() string {
	return "uniq_by_line"
}

func (p *UniqByLine) Process(issues []result.Issue) ([]result.Issue, error) {
	if !p.cfg.Output.UniqByLine {
		return issues, nil
	}

	return filterIssuesUnsafe(issues, p.shouldPassIssue), nil
}

func (*UniqByLine) Finish() {}

func (p *UniqByLine) shouldPassIssue(issue *result.Issue) bool {
	if issue.Replacement != nil && p.cfg.Issues.NeedFix {
		// if issue will be auto-fixed we shouldn't collapse issues:
		// e.g. one line can contain 2 misspellings, they will be in 2 issues and misspell should fix both of them.
		return true
	}

	if p.fileLineCounter.GetCount(issue) == uniqByLineLimit {
		return false
	}

	p.fileLineCounter.Increment(issue)

	return true
}

type fileLineCounter map[string]map[int]int

func (f fileLineCounter) GetCount(issue *result.Issue) int {
	return f.getCounter(issue)[issue.Line()]
}

func (f fileLineCounter) Increment(issue *result.Issue) {
	f.getCounter(issue)[issue.Line()]++
}

func (f fileLineCounter) getCounter(issue *result.Issue) map[int]int {
	lc := f[issue.FilePath()]

	if lc == nil {
		lc = map[int]int{}
		f[issue.FilePath()] = lc
	}

	return lc
}
