package processors

import (
	"fmt"
	"path/filepath"

	"github.com/golangci/golangci-lint/pkg/fsutils"
	"github.com/golangci/golangci-lint/pkg/result"
)

type PathPrettifier struct {
	root string
}

var _ Processor = PathPrettifier{}

func NewPathPrettifier() *PathPrettifier {
	root, err := fsutils.Getwd()
	if err != nil {
		panic(fmt.Sprintf("Can't get working dir: %s", err))
	}
	return &PathPrettifier{
		root: root,
	}
}

func (p PathPrettifier) Name() string {
	return "path_prettifier"
}

func (p PathPrettifier) Process(issues []result.Issue) ([]result.Issue, error) {
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

func (p PathPrettifier) Finish() {}
