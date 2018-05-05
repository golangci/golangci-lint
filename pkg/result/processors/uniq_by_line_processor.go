package processors

import (
	"fmt"

	"github.com/golangci/golangci-lint/pkg/result"
)

type UniqByLineProcessor struct{}

var _ Processor = UniqByLineProcessor{}

func (p UniqByLineProcessor) Name() string {
	return "uniq_by_line"
}

func (p UniqByLineProcessor) Process(results []result.Result) ([]result.Result, error) {
	fli := makeFilesToLinesToIssuesMap(results)

	retResults := []result.Result{}
	for _, res := range results {
		newRes := res
		newRes.Issues = []result.Issue{}
		for _, i := range res.Issues {
			lineIssues := fli[i.File][i.LineNumber]
			if len(lineIssues) == 0 {
				return nil, fmt.Errorf("bug in by line uniqalization")
			}

			if lineIssues[0] == i { // Use first issue for line
				newRes.Issues = append(newRes.Issues, i)
			}
		}
		retResults = append(retResults, newRes)
	}

	return retResults, nil
}
