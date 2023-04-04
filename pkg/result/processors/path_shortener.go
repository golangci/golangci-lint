package processors

import (
	"fmt"
	"strings"

	"github.com/golangci/golangci-lint/pkg/fsutils"
	"github.com/golangci/golangci-lint/pkg/result"
)

type PathShortener struct {
	wd string
}

var _ Processor = PathShortener{}

func NewPathShortener() *PathShortener {
	wd, err := fsutils.Getwd()
	if err != nil {
		panic(fmt.Sprintf("Can't get working dir: %s", err))
	}
	return &PathShortener{
		wd: wd,
	}
}

func (p PathShortener) Name() string {
	return "path_shortener"
}

func (p PathShortener) Process(issues []result.Issue) ([]result.Issue, error) {
	return transformIssues(issues, func(i *result.Issue) *result.Issue {
		newI := i
		newI.Text = strings.ReplaceAll(newI.Text, p.wd+"/", "")
		newI.Text = strings.ReplaceAll(newI.Text, p.wd, "")
		return newI
	}), nil
}

func (p PathShortener) Finish() {}
