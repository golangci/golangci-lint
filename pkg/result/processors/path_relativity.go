package processors

import (
	"fmt"
	"path/filepath"

	"github.com/golangci/golangci-lint/pkg/fsutils"
	"github.com/golangci/golangci-lint/pkg/logutils"
	"github.com/golangci/golangci-lint/pkg/result"
)

var _ Processor = (*PathRelativity)(nil)

// PathRelativity computes [result.Issue.RelativePath] and  [result.Issue.WorkingDirectoryRelativePath],
// based on the base path.
type PathRelativity struct {
	log      logutils.Log
	basePath string
	wd       string
}

func NewPathRelativity(log logutils.Log, basePath string) (*PathRelativity, error) {
	wd, err := fsutils.Getwd()
	if err != nil {
		return nil, fmt.Errorf("error getting working directory: %w", err)
	}

	return &PathRelativity{
		log:      log.Child(logutils.DebugKeyPathRelativity),
		basePath: basePath,
		wd:       wd,
	}, nil
}

func (p *PathRelativity) Process(issues []result.Issue) ([]result.Issue, error) {
	return transformIssues(issues, func(issue *result.Issue) *result.Issue {
		newIssue := issue

		var err error
		newIssue.RelativePath, err = filepath.Rel(p.basePath, issue.FilePath())
		if err != nil {
			p.log.Warnf("relative path (basepath): %v", err)
			return nil
		}

		newIssue.WorkingDirectoryRelativePath, err = filepath.Rel(p.wd, issue.FilePath())
		if err != nil {
			p.log.Warnf("relative path (wd): %v", err)
			return nil
		}

		return newIssue
	}), nil
}

func (*PathRelativity) Name() string {
	return "path_relativity"
}

func (*PathRelativity) Finish() {}
