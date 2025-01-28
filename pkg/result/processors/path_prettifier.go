package processors

import (
	"github.com/golangci/golangci-lint/pkg/fsutils"
	"github.com/golangci/golangci-lint/pkg/logutils"
	"github.com/golangci/golangci-lint/pkg/result"
)

var _ Processor = (*PathPrettifier)(nil)

// PathPrettifier modifies report file path to be relative to the base path.
// Also handles the `output.path-prefix` option.
type PathPrettifier struct {
	prefix string
	log    logutils.Log
}

func NewPathPrettifier(log logutils.Log, prefix string) *PathPrettifier {
	return &PathPrettifier{
		prefix: prefix,
		log:    log.Child(logutils.DebugKeyPathPrettifier),
	}
}

func (*PathPrettifier) Name() string {
	return "path_prettifier"
}

func (p *PathPrettifier) Process(issues []result.Issue) ([]result.Issue, error) {
	return transformIssues(issues, func(issue *result.Issue) *result.Issue {
		newIssue := issue

		newIssue.Pos.Filename = fsutils.WithPathPrefix(p.prefix, issue.RelativePath)

		return newIssue
	}), nil
}

func (*PathPrettifier) Finish() {}
