package golinters

import (
	"testing"

	"github.com/golangci/golangci-lint/pkg/result"
)

func TestGovetSimple(t *testing.T) {
	const source = `package p

import "os"

func f() error {
  return &os.PathError{"first", "path", os.ErrNotExist}
}
`

	ExpectIssues(t, govet, source, []result.Issue{
		NewIssue("govet", "os.PathError composite literal uses unkeyed fields", 6),
	})
}
