package processors

import (
	"github.com/golangci/golangci-lint/pkg/logutils"
	"github.com/golangci/golangci-lint/pkg/result"
)

var _ Processor = (*PathPrettifier)(nil)

// PathPrettifier modifies report file path with the shortest relative path.
type PathPrettifier struct {
	log logutils.Log
}

func NewPathPrettifier(log logutils.Log) *PathPrettifier {
	return &PathPrettifier{log: log.Child(logutils.DebugKeyPathPrettifier)}
}

func (*PathPrettifier) Name() string {
	return "path_prettifier"
}

func (*PathPrettifier) Process(issues []result.Issue) ([]result.Issue, error) {
	return transformIssues(issues, func(issue *result.Issue) *result.Issue {
		newIssue := issue
		newIssue.Pos.Filename = issue.RelativePath

		return newIssue
	}), nil
}

func (*PathPrettifier) Finish() {}
