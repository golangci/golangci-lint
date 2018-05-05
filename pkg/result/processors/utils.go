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
