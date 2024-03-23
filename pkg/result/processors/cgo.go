package processors

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/golangci/golangci-lint/pkg/goutil"
	"github.com/golangci/golangci-lint/pkg/result"
)

var _ Processor = (*Cgo)(nil)

type Cgo struct {
	goCacheDir string
}

func NewCgo(goenv *goutil.Env) *Cgo {
	return &Cgo{
		goCacheDir: goenv.Get(goutil.EnvGoCache),
	}
}

func (Cgo) Name() string {
	return "cgo"
}

func (p Cgo) Process(issues []result.Issue) ([]result.Issue, error) {
	return filterIssuesErr(issues, p.shouldPassIssue)
}

func (Cgo) Finish() {}

func (p Cgo) shouldPassIssue(issue *result.Issue) (bool, error) {
	// some linters (e.g. gosec, deadcode) return incorrect filepaths for cgo issues,
	// also cgo files have strange issues looking like false positives.

	// cache dir contains all preprocessed files including cgo files

	issueFilePath := issue.FilePath()
	if !filepath.IsAbs(issue.FilePath()) {
		absPath, err := filepath.Abs(issue.FilePath())
		if err != nil {
			return false, fmt.Errorf("failed to build abs path for %q: %w", issue.FilePath(), err)
		}
		issueFilePath = absPath
	}

	if p.goCacheDir != "" && strings.HasPrefix(issueFilePath, p.goCacheDir) {
		return false, nil
	}

	if filepath.Base(issue.FilePath()) == "_cgo_gotypes.go" {
		// skip cgo warning for go1.10
		return false, nil
	}

	return true, nil
}
