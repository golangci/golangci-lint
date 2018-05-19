package printers

import "github.com/golangci/golangci-lint/pkg/result"

type Printer interface {
	Print(issues <-chan result.Issue) (bool, error)
}
