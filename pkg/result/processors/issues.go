package processors

import (
	"fmt"

	"github.com/golangci/golangci-lint/pkg/result"
)

func filterIssues(issues []result.Issue, filter func(i *result.Issue) bool) []result.Issue {
	retIssues := make([]result.Issue, 0, len(issues))
	for i := range issues {
		if filter(&issues[i]) {
			retIssues = append(retIssues, issues[i])
		}
	}

	return retIssues
}

func filterIssuesErr(issues []result.Issue, filter func(i *result.Issue) (bool, error)) ([]result.Issue, error) {
	retIssues := make([]result.Issue, 0, len(issues))
	for i := range issues {
		ok, err := filter(&issues[i])
		if err != nil {
			return nil, fmt.Errorf("can't filter issue %#v: %w", issues[i], err)
		}

		if ok {
			retIssues = append(retIssues, issues[i])
		}
	}

	return retIssues, nil
}

func transformIssues(issues []result.Issue, transform func(i *result.Issue) *result.Issue) []result.Issue {
	retIssues := make([]result.Issue, 0, len(issues))
	for i := range issues {
		newI := transform(&issues[i])
		if newI != nil {
			retIssues = append(retIssues, *newI)
		}
	}

	return retIssues
}
