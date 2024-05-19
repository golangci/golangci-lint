package printers

import (
	"fmt"
	"io"
	"path/filepath"

	"github.com/golangci/golangci-lint/pkg/result"
)

const defaultGithubSeverity = "error"

type GitHubAction struct {
	w io.Writer
}

// NewGitHubAction output format outputs issues according to GitHub actions.
// Deprecated
func NewGitHubAction(w io.Writer) *GitHubAction {
	return &GitHubAction{w: w}
}

func (p *GitHubAction) Print(issues []result.Issue) error {
	for ind := range issues {
		_, err := fmt.Fprintln(p.w, formatIssueAsGitHub(&issues[ind]))
		if err != nil {
			return err
		}
	}
	return nil
}

// print each line as: ::error file=app.js,line=10,col=15::Something went wrong
func formatIssueAsGitHub(issue *result.Issue) string {
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
