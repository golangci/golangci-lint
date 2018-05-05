package golinters

import (
	"testing"

	"github.com/golangci/golangci-lint/pkg/result"
)

func TestGofmtIssueFound(t *testing.T) {
	const source = `package p

func noFmt() error {
return nil
}
`

	ExpectIssues(t, gofmt{}, source, []result.Issue{NewIssue("gofmt", "File is not gofmt-ed with -s", 4)})
}

func TestGofmtNoIssue(t *testing.T) {
	const source = `package p

func fmted() error {
	return nil
}
`

	ExpectIssues(t, gofmt{}, source, []result.Issue{})
}

func TestGoimportsIssueFound(t *testing.T) {
	const source = `package p
func noFmt() error {return nil}
`

	lint := gofmt{useGoimports: true}
	ExpectIssues(t, lint, source, []result.Issue{NewIssue("goimports", "File is not goimports-ed", 2)})
}

func TestGoimportsNoIssue(t *testing.T) {
	const source = `package p

func fmted() error {
	return nil
}
`

	lint := gofmt{useGoimports: true}
	ExpectIssues(t, lint, source, []result.Issue{})
}
