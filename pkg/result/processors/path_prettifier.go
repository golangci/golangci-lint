package processors

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/golangci/golangci-lint/pkg/result"
)

type PathPrettifier struct {
	root string
}

var _ Processor = PathPrettifier{}

func NewPathPrettifier() *PathPrettifier {
	root, err := os.Getwd()
	if err != nil {
		panic(fmt.Errorf("Can't get working dir: %s", err))
	}
	return &PathPrettifier{
		root: root,
	}
}

func (p PathPrettifier) Name() string {
	return "path_prettifier"
}

func (p PathPrettifier) Process(issues []result.Issue) ([]result.Issue, error) {
	return transformIssues(issues, func(i *result.Issue) *result.Issue {
		if !filepath.IsAbs(i.File) {
			return i
		}

		rel, err := filepath.Rel(p.root, i.File)
		if err != nil {
			return i
		}

		newI := i
		newI.File = rel
		return newI
	}), nil
}
