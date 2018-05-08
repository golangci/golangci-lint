package processors

import (
	"github.com/golangci/golangci-lint/pkg/result"
)

type lineToCount map[int]int
type fileToLineToCount map[string]lineToCount

type UniqByLine struct {
	flc fileToLineToCount
}

func NewUniqByLine() *UniqByLine {
	return &UniqByLine{
		flc: fileToLineToCount{},
	}
}

var _ Processor = &UniqByLine{}

func (p UniqByLine) Name() string {
	return "uniq_by_line"
}

func (p *UniqByLine) Process(issues []result.Issue) ([]result.Issue, error) {
	return filterIssues(issues, func(i *result.Issue) bool {
		lc := p.flc[i.File]
		if lc == nil {
			lc = lineToCount{}
			p.flc[i.File] = lc
		}

		const limit = 1
		count := lc[i.LineNumber]
		if count == limit {
			return false
		}

		lc[i.LineNumber]++
		return true
	}), nil
}

func (p UniqByLine) Finish() {}
