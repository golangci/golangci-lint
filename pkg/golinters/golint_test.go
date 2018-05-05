package golinters

import (
	"testing"

	"github.com/golangci/golangci-lint/pkg/result"
)

func TestGolintSimple(t *testing.T) {
	const source = `package p
	var v_1 string`

	ExpectIssues(t, golint, source,
		[]result.Issue{NewIssue("golint", "don't use underscores in Go names; var v_1 should be v1", 2)})
}
