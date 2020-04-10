package printers

import (
	"context"
	"fmt"
	"github.com/golangci/golangci-lint/pkg/logutils"
	"github.com/golangci/golangci-lint/pkg/result"
)

type github struct {
}

// Github output format outputs issues according to Github actions format:
// https://help.github.com/en/actions/reference/workflow-commands-for-github-actions#setting-an-error-message
func NewGithub() Printer {
	return &github{}
}

// print each line as: ::error file=app.js,line=10,col=15::Something went wrong
func formatIssueAsGithub(issue result.Issue) string {
	result := fmt.Sprintf("::error file=%s,line=%d", issue.FilePath(), issue.Line())
	if issue.Pos.Column != 0 {
		result += fmt.Sprintf(",col=%d", issue.Pos.Column)
	}

	result += fmt.Sprintf("::%s (%s)", issue.Text, issue.FromLinter)
	return result
}

func (g *github) Print(ctx context.Context, issues []result.Issue) error {
	for _, issue := range issues {
		_, err := fmt.Fprintln(logutils.StdOut, formatIssueAsGithub(issue))
		if err != nil {
			return err
		}
	}
	return nil
}
