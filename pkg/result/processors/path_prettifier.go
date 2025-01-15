package processors

import (
	"path/filepath"

	"github.com/golangci/golangci-lint/pkg/fsutils"
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

func (p *PathPrettifier) Process(issues []result.Issue) ([]result.Issue, error) {
	return transformIssues(issues, func(issue *result.Issue) *result.Issue {
		if !filepath.IsAbs(issue.FilePath()) {
			return issue
		}

		rel, err := fsutils.ShortestRelPath(issue.FilePath(), "")
		if err != nil {
			p.log.Warnf("shortest relative path: %v", err)
			return issue
		}

		newIssue := issue
		newIssue.Pos.Filename = rel
		return newIssue
	}), nil
}

func (*PathPrettifier) Finish() {}
