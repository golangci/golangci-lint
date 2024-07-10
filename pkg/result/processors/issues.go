package processors

import (
	"fmt"

	"github.com/golangci/golangci-lint/pkg/result"
)

func filterIssues(issues []result.Issue, filter func(issue *result.Issue) bool) []result.Issue {
	retIssues := make([]result.Issue, 0, len(issues))
	for i := range issues {
		if issues[i].FromLinter == typeCheckName {
			// don't hide typechecking errors in generated files: users expect to see why the project isn't compiling
			retIssues = append(retIssues, issues[i])
			continue
		}

		if filter(&issues[i]) {
			retIssues = append(retIssues, issues[i])
		}
	}

	return retIssues
}

func filterIssuesUnsafe(issues []result.Issue, filter func(issue *result.Issue) bool) []result.Issue {
	retIssues := make([]result.Issue, 0, len(issues))
	for i := range issues {
		if filter(&issues[i]) {
			retIssues = append(retIssues, issues[i])
		}
	}

	return retIssues
}

func filterIssuesErr(issues []result.Issue, filter func(issue *result.Issue) (bool, error)) ([]result.Issue, error) {
	retIssues := make([]result.Issue, 0, len(issues))
	for i := range issues {
		if issues[i].FromLinter == typeCheckName {
			// don't hide typechecking errors in generated files: users expect to see why the project isn't compiling
			retIssues = append(retIssues, issues[i])
			continue
		}

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

func transformIssues(issues []result.Issue, transform func(issue *result.Issue) *result.Issue) []result.Issue {
	retIssues := make([]result.Issue, 0, len(issues))
	for i := range issues {
		newIssue := transform(&issues[i])
		if newIssue != nil {
			retIssues = append(retIssues, *newIssue)
		}
	}

	return retIssues
}
