package processors

import (
	"github.com/golangci/golangci-lint/pkg/goutils"
	"github.com/golangci/golangci-lint/pkg/result"
)

type Cgo struct {
}

var _ Processor = Cgo{}

func NewCgo() *Cgo {
	return &Cgo{}
}

func (p Cgo) Name() string {
	return "cgo"
}

func (p Cgo) Process(issues []result.Issue) ([]result.Issue, error) {
	return filterIssues(issues, func(i *result.Issue) bool {
		// some linters (.e.g gas, deadcode) return incorrect filepaths for cgo issues,
		// it breaks next processing, so skip them
		return !goutils.IsCgoFilename(i.FilePath())
	}), nil
}

func (Cgo) Finish() {}
