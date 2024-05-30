package processors

import (
	"path/filepath"

	"github.com/golangci/golangci-lint/pkg/fsutils"
	"github.com/golangci/golangci-lint/pkg/result"
)

var _ Processor = (*PathPrettifier)(nil)

type PathPrettifier struct {
}

func NewPathPrettifier() *PathPrettifier {
	return &PathPrettifier{}
}

func (PathPrettifier) Name() string {
	return "path_prettifier"
}

func (PathPrettifier) Process(issues []result.Issue) ([]result.Issue, error) {
	return transformIssues(issues, func(issue *result.Issue) *result.Issue {
		if !filepath.IsAbs(issue.FilePath()) {
			return issue
		}

		rel, err := fsutils.ShortestRelPath(issue.FilePath(), "")
		if err != nil {
			return issue
		}

		newIssue := issue
		newIssue.Pos.Filename = rel
		return newIssue
	}), nil
}

func (PathPrettifier) Finish() {}
