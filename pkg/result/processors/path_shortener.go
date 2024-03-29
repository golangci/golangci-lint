package processors

import (
	"fmt"
	"strings"

	"github.com/golangci/golangci-lint/pkg/fsutils"
	"github.com/golangci/golangci-lint/pkg/result"
)

var _ Processor = (*PathShortener)(nil)

type PathShortener struct {
	wd string
}

func NewPathShortener() *PathShortener {
	wd, err := fsutils.Getwd()
	if err != nil {
		panic(fmt.Sprintf("Can't get working dir: %s", err))
	}

	return &PathShortener{wd: wd}
}

func (PathShortener) Name() string {
	return "path_shortener"
}

func (p PathShortener) Process(issues []result.Issue) ([]result.Issue, error) {
	return transformIssues(issues, func(issue *result.Issue) *result.Issue {
		newIssue := issue
		newIssue.Text = strings.ReplaceAll(newIssue.Text, p.wd+"/", "")
		newIssue.Text = strings.ReplaceAll(newIssue.Text, p.wd, "")
		return newIssue
	}), nil
}

func (PathShortener) Finish() {}
