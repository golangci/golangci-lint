package processors

import (
	"fmt"
	"path/filepath"

	"github.com/golangci/golangci-lint/v2/pkg/fsutils"
	"github.com/golangci/golangci-lint/v2/pkg/logutils"
	"github.com/golangci/golangci-lint/v2/pkg/result"
)

var _ Processor = (*PathRelativity)(nil)

// PathRelativity computes [result.Issue.RelativePath] and [result.Issue.WorkingDirectoryRelativePath],
// based on the base path.
type PathRelativity struct {
	log              logutils.Log
	basePath         string
	workingDirectory string
}

func NewPathRelativity(log logutils.Log, basePath string) (*PathRelativity, error) {
	wd, err := fsutils.Getwd()
	if err != nil {
		return nil, fmt.Errorf("error getting working directory: %w", err)
	}

	return &PathRelativity{
		log:              log.Child(logutils.DebugKeyPathRelativity),
		basePath:         evalSymlinks(basePath),
		workingDirectory: wd,
	}, nil
}

func (*PathRelativity) Name() string {
	return "path_relativity"
}

func (p *PathRelativity) Process(issues []*result.Issue) ([]*result.Issue, error) {
	return transformIssues(issues, func(issue *result.Issue) *result.Issue {
		newIssue := *issue

		filePath := evalSymlinks(issue.FilePath())

		var err error
		newIssue.RelativePath, err = filepath.Rel(p.basePath, filePath)
		if err != nil {
			p.log.Warnf("Getting relative path (basepath): %v", err)
			return nil
		}

		newIssue.WorkingDirectoryRelativePath, err = filepath.Rel(p.workingDirectory, filePath)
		if err != nil {
			p.log.Warnf("Getting relative path (wd): %v", err)
			return nil
		}

		return &newIssue
	}), nil
}

func (*PathRelativity) Finish() {}

// evalSymlinks resolves symlinks on a best-effort basis.
//
// The base path and the issue paths can reach the same directory through different spellings:
// [fsutils.Getwd] resolves symlinks, and `git rev-parse --show-toplevel` returns a resolved path too,
// while the file names reported by the linters keep the spelling used by the Go toolchain.
// Resolving both sides prevents [filepath.Rel] from comparing two spellings of the same directory,
// which produced paths leaving the base path and coming back into it.
//
// Paths that cannot be resolved are returned unchanged: they may simply not exist on disk,
// and that is not a reason to drop the issue.
func evalSymlinks(path string) string {
	if path == "" {
		return path
	}

	resolved, err := fsutils.EvalSymlinks(path)
	if err != nil {
		return path
	}

	return resolved
}
