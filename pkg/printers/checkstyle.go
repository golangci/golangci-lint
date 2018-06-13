package printers

import (
	"context"

	"github.com/golangci/golangci-lint/pkg/result"
)

type Checkstyle struct{}

func NewCheckstyle() *Checkstyle {
	return &Checkstyle{}
}

func (Checkstyle) Print(ctx context.Context, issues <-chan result.Issue) (bool, error) {
	return false, nil
}
