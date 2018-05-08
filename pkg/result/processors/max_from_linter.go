package processors

import (
	"github.com/golangci/golangci-lint/pkg/result"
	"github.com/sirupsen/logrus"
)

type MaxFromLinter struct {
	lc    linterToCountMap
	limit int
}

var _ Processor = &MaxFromLinter{}

func NewMaxFromLinter(limit int) *MaxFromLinter {
	return &MaxFromLinter{
		lc:    linterToCountMap{},
		limit: limit,
	}
}

func (p MaxFromLinter) Name() string {
	return "max_from_linter"
}

func (p *MaxFromLinter) Process(issues []result.Issue) ([]result.Issue, error) {
	if p.limit <= 0 { // no limit
		return issues, nil
	}

	return filterIssues(issues, func(i *result.Issue) bool {
		p.lc[i.FromLinter]++ // always inc for stat
		return p.lc[i.FromLinter] <= p.limit
	}), nil
}

func (p MaxFromLinter) Finish() {
	for linter, count := range p.lc {
		if count > p.limit {
			logrus.Infof("%d/%d issues from linter %s were hidden, use --max-issues-per-linter",
				count-p.limit, count, linter)
		}
	}
}
