package processors

import (
	"github.com/golangci/golangci-lint/pkg/result"
	"github.com/sirupsen/logrus"
)

type textToCountMap map[string]int

type MaxSameIssues struct {
	tc    textToCountMap
	limit int
}

var _ Processor = &MaxSameIssues{}

func NewMaxSameIssues(limit int) *MaxSameIssues {
	return &MaxSameIssues{
		tc:    textToCountMap{},
		limit: limit,
	}
}

func (MaxSameIssues) Name() string {
	return "max_same_issues"
}

func (p *MaxSameIssues) Process(issues []result.Issue) ([]result.Issue, error) {
	if p.limit <= 0 { // no limit
		return issues, nil
	}

	return filterIssues(issues, func(i *result.Issue) bool {
		p.tc[i.Text]++ // always inc for stat
		return p.tc[i.Text] <= p.limit
	}), nil
}

func (p MaxSameIssues) Finish() {
	for text, count := range p.tc {
		if count > p.limit {
			logrus.Infof("%d/%d issues with text %q were hidden, use --max-same-issues",
				count-p.limit, count, text)
		}
	}
}
