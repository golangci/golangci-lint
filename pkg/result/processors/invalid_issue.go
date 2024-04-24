package processors

import (
	"path/filepath"

	"github.com/golangci/golangci-lint/pkg/logutils"
	"github.com/golangci/golangci-lint/pkg/result"
)

var _ Processor = (*InvalidIssue)(nil)

type InvalidIssue struct {
	log logutils.Log
}

func NewInvalidIssue(log logutils.Log) *InvalidIssue {
	return &InvalidIssue{log: log}
}

func (InvalidIssue) Name() string {
	return "invalid_issue"
}

func (p InvalidIssue) Process(issues []result.Issue) ([]result.Issue, error) {
	tcIssues := filterIssues(issues, func(issue *result.Issue) bool {
		return issue.FromLinter == typeCheckName
	})

	if len(tcIssues) > 0 {
		return tcIssues, nil
	}

	return filterIssuesErr(issues, p.shouldPassIssue)
}

func (InvalidIssue) Finish() {}

func (p InvalidIssue) shouldPassIssue(issue *result.Issue) (bool, error) {
	if issue.FilePath() == "" {
		p.log.Warnf("no file path for the issue: probably a bug inside the linter %q: %#v", issue.FromLinter, issue)

		return false, nil
	}

	if filepath.Base(issue.FilePath()) == "go.mod" {
		return true, nil
	}

	if !isGoFile(issue.FilePath()) {
		p.log.Infof("issue related to file %s is skipped", issue.FilePath())

		return false, nil
	}

	return true, nil
}

func isGoFile(name string) bool {
	return filepath.Ext(name) == ".go"
}
