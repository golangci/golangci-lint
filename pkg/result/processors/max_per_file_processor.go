package processors

import "github.com/golangci/golangci-lint/pkg/result"

type MaxLinterIssuesPerFile struct{}

var _ Processor = MaxLinterIssuesPerFile{}

type fileToIssuesMap map[string][]result.Issue

func (p MaxLinterIssuesPerFile) Name() string {
	return "max_issues_per_file"
}

func (p MaxLinterIssuesPerFile) makeFileToIssuesMap(res result.Result) fileToIssuesMap {
	fti := fileToIssuesMap{}
	for _, i := range res.Issues {
		fti[i.File] = append(fti[i.File], i)
	}

	return fti
}

func (p MaxLinterIssuesPerFile) processResult(res result.Result) result.Result {
	if len(res.Issues) == 0 {
		return res
	}

	if res.MaxIssuesPerFile == 0 {
		return res // Nothing to process
	}

	fti := p.makeFileToIssuesMap(res)
	for file, fileIssues := range fti {
		if len(fileIssues) > res.MaxIssuesPerFile {
			fti[file] = fileIssues[:res.MaxIssuesPerFile]
		}
	}

	filteredIssues := []result.Issue{}
	for _, issues := range fti {
		filteredIssues = append(filteredIssues, issues...)
	}

	res.Issues = filteredIssues
	return res
}

func (p MaxLinterIssuesPerFile) Process(results []result.Result) ([]result.Result, error) {
	newResults := []result.Result{}

	for _, res := range results {
		newResults = append(newResults, p.processResult(res))
	}

	return newResults, nil
}
