package printers

import (
	"fmt"
	"io"
	"path/filepath"

	"github.com/golangci/golangci-lint/pkg/result"
)

type github struct {
	w io.Writer
}

const defaultGithubSeverity = "error"

// NewGithub output format outputs issues according to GitHub actions format:
// https://docs.github.com/en/actions/using-workflows/workflow-commands-for-github-actions#setting-an-error-message
func NewGithub(w io.Writer) Printer {
	return &github{w: w}
}

// print each line as: ::error file=app.js,line=10,col=15::Something went wrong
func formatIssueAsGithub(issue *result.Issue) string {
	severity := defaultGithubSeverity
	if issue.Severity != "" {
		severity = issue.Severity
	}

	// Convert backslashes to forward slashes.
	// This is needed when running on windows.
	// Otherwise, GitHub won't be able to show the annotations pointing to the file path with backslashes.
	file := filepath.ToSlash(issue.FilePath())

	ret := fmt.Sprintf("::%s file=%s,line=%d", severity, file, issue.Line())
	if issue.Pos.Column != 0 {
		ret += fmt.Sprintf(",col=%d", issue.Pos.Column)
	}

	ret += fmt.Sprintf("::%s (%s)", issue.Text, issue.FromLinter)
	return ret
}

func (p *github) Print(issues []result.Issue) error {
	for ind := range issues {
		_, err := fmt.Fprintln(p.w, formatIssueAsGithub(&issues[ind]))
		if err != nil {
			return err
		}
	}
	return nil
}
