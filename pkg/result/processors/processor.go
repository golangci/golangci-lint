package processors

import "github.com/golangci/golangci-lint/pkg/result"

type Processor interface {
	Process(results []result.Result) ([]result.Result, error)
	Name() string
}
