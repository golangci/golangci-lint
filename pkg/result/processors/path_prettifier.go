package processors

import (
	"regexp"

	"github.com/golangci/golangci-lint/pkg/result"
)

type ExcludeProcessor struct {
	pattern *regexp.Regexp
}

var _ Processor = ExcludeProcessor{}

func NewExcludeProcessor(pattern string) *ExcludeProcessor {
	return &ExcludeProcessor{
		pattern: regexp.MustCompile(pattern),
	}
}

func (p ExcludeProcessor) Name() string {
	return "exclude"
}

func (p ExcludeProcessor) processResult(res result.Result) result.Result {
	newRes := res
	newRes.Issues = []result.Issue{}
	for _, i := range res.Issues {
		if !p.pattern.MatchString(i.Text) {
			newRes.Issues = append(newRes.Issues, i)
		}
	}

	return newRes
}

func (p ExcludeProcessor) Process(results []result.Result) ([]result.Result, error) {
	retResults := []result.Result{}
	for _, res := range results {
		retResults = append(retResults, p.processResult(res))
	}

	return retResults, nil
}
