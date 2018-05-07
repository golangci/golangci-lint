package processors

import "github.com/golangci/golangci-lint/pkg/result"

type linesToIssuesMap map[int][]result.Issue
type filesToLinesToIssuesMap map[string]linesToIssuesMap

func makeFilesToLinesToIssuesMap(results []result.Result) filesToLinesToIssuesMap {
	fli := filesToLinesToIssuesMap{}
	for _, res := range results {
		for _, i := range res.Issues {
			if fli[i.File] == nil {
				fli[i.File] = linesToIssuesMap{}
			}
			li := fli[i.File]
			li[i.LineNumber] = append(li[i.LineNumber], i)
		}
	}
	return fli
}

func filterIssues(issues []result.Issue, filter func(i *result.Issue) bool) []result.Issue {
	retIssues := make([]result.Issue, 0, len(issues))
	for _, i := range issues {
		if filter(&i) {
			retIssues = append(retIssues, i)
		}
	}

	return retIssues
}

func filterIssuesErr(issues []result.Issue, filter func(i *result.Issue) (bool, error)) ([]result.Issue, error) {
	retIssues := make([]result.Issue, 0, len(issues))
	for _, i := range issues {
		ok, err := filter(&i)
		if err != nil {
			return nil, err
		}
		if ok {
			retIssues = append(retIssues, i)
		}
	}

	return retIssues, nil
}

func transformIssues(issues []result.Issue, transform func(i *result.Issue) *result.Issue) []result.Issue {
	retIssues := make([]result.Issue, 0, len(issues))
	for _, i := range issues {
		newI := transform(&i)
		if newI != nil {
			retIssues = append(retIssues, *newI)
		}
	}

	return retIssues
}
