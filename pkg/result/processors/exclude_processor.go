package processors

import (
	"log"
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
		log.Fatalf("Can't get working dir: %s", err)
	}
	return &PathPrettifier{
		root: root,
	}
}

func (p PathPrettifier) Name() string {
	return "path_prettifier"
}

func (p PathPrettifier) processResult(res result.Result) result.Result {
	newRes := res
	newRes.Issues = []result.Issue{}
	for _, i := range res.Issues {
		if filepath.IsAbs(i.File) {
			if rel, err := filepath.Rel(p.root, i.File); err == nil {
				newI := i
				newI.File = rel
				newRes.Issues = append(newRes.Issues, newI)
				continue
			}
		}

		newRes.Issues = append(newRes.Issues, i)
	}

	return newRes
}

func (p PathPrettifier) Process(results []result.Result) ([]result.Result, error) {
	retResults := []result.Result{}
	for _, res := range results {
		retResults = append(retResults, p.processResult(res))
	}

	return retResults, nil
}
