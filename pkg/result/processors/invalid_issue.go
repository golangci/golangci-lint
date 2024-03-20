package processors

import (
	"path/filepath"

	"github.com/golangci/golangci-lint/pkg/logutils"
	"github.com/golangci/golangci-lint/pkg/result"
)

var _ Processor = InvalidIssue{}

type InvalidIssue struct {
	log logutils.Log
}

func NewInvalidIssue(log logutils.Log) *InvalidIssue {
	return &InvalidIssue{log: log}
}

func (p InvalidIssue) Process(issues []result.Issue) ([]result.Issue, error) {
	return filterIssuesErr(issues, p.shouldPassIssue)
}

func (p InvalidIssue) Name() string {
	return "invalid_issue"
}

func (p InvalidIssue) Finish() {}

func (p InvalidIssue) shouldPassIssue(issue *result.Issue) (bool, error) {
	if issue.FromLinter == "typecheck" {
		return true, nil
	}

	if issue.FilePath() == "" {
		// contextcheck has a known bug https://github.com/kkHAIKE/contextcheck/issues/21
		if issue.FromLinter != "contextcheck" {
			p.log.Warnf("no file path for the issue: probably a bug inside the linter %q: %#v", issue.FromLinter, issue)
		}

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
