package processors

import (
	"path/filepath"
	"strings"

	"github.com/golangci/golangci-lint/pkg/goutil"
	"github.com/golangci/golangci-lint/pkg/result"
)

var _ Processor = (*Cgo)(nil)

// Cgo some linters (e.g. gosec, deadcode) return incorrect filepaths for cgo issues,
// also cgo files have strange issues looking like false positives.
//
// Require absolute filepath.
type Cgo struct {
	goCacheDir string
}

func NewCgo(goenv *goutil.Env) *Cgo {
	return &Cgo{
		goCacheDir: goenv.Get(goutil.EnvGoCache),
	}
}

func (*Cgo) Name() string {
	return "cgo"
}

func (p *Cgo) Process(issues []result.Issue) ([]result.Issue, error) {
	return filterIssuesErr(issues, p.shouldPassIssue)
}

func (*Cgo) Finish() {}

func (p *Cgo) shouldPassIssue(issue *result.Issue) (bool, error) {
	// [p.goCacheDir] contains all preprocessed files including cgo files.
	if p.goCacheDir != "" && strings.HasPrefix(issue.FilePath(), p.goCacheDir) {
		return false, nil
	}

	if filepath.Base(issue.FilePath()) == "_cgo_gotypes.go" {
		// skip cgo warning for go1.10
		return false, nil
	}

	return true, nil
}
