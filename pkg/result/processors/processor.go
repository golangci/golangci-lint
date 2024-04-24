package processors

import (
	"github.com/golangci/golangci-lint/pkg/result"
)

const typeCheckName = "typecheck"

type Processor interface {
	Process(issues []result.Issue) ([]result.Issue, error)
	Name() string
	Finish()
}
