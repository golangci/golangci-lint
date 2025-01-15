package processors

import (
	"github.com/golangci/golangci-lint/pkg/fsutils"
	"github.com/golangci/golangci-lint/pkg/result"
)

var _ Processor = (*PathPrefixer)(nil)

// PathPrefixer adds a customizable path prefix to report file paths for user facing.
// It uses the shortest relative paths and `path-prefix` option.
type PathPrefixer struct {
	prefix string
}

// NewPathPrefixer returns a new path prefixer for the provided string
func NewPathPrefixer(prefix string) *PathPrefixer {
	return &PathPrefixer{prefix: prefix}
}

// Name returns the name of this processor
func (*PathPrefixer) Name() string {
	return "path_prefixer"
}

// Process adds the prefix to each path
func (p *PathPrefixer) Process(issues []result.Issue) ([]result.Issue, error) {
	if p.prefix == "" {
		return issues, nil
	}

	for i := range issues {
		issues[i].Pos.Filename = fsutils.WithPathPrefix(p.prefix, issues[i].Pos.Filename)
	}

	return issues, nil
}

// Finish is implemented to satisfy the Processor interface
func (*PathPrefixer) Finish() {}
