package main

import (
	"cmp"
	"slices"
	"strconv"
	"strings"

	"github.com/securego/gosec/v2/issue"
)

// handle ranges
func extractLineNumber(s string) int {
	lineNumber, _ := strconv.Atoi(strings.Split(s, "-")[0])
	return lineNumber
}

// sortIssues sorts the issues by severity in descending order
func sortIssues(issues []*issue.Issue) {
	slices.SortFunc(issues, func(i, j *issue.Issue) int {
		return -cmp.Or(
			cmp.Compare(i.Severity, j.Severity),
			cmp.Compare(i.What, j.What),
			cmp.Compare(i.File, j.File),
			cmp.Compare(extractLineNumber(i.Line), extractLineNumber(j.Line)),
		)
	})
}
