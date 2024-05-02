package printers

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/golangci/golangci-lint/pkg/result"
)

const defaultGitHubSeverity = "error"

const filenameGitHubActionProblemMatchers = "golangci-lint-action-problem-matchers.json"

// GitHubProblemMatchers defines the root of problem matchers.
// - https://github.com/actions/toolkit/blob/main/docs/problem-matchers.md
// - https://github.com/actions/toolkit/blob/main/docs/commands.md#problem-matchers
type GitHubProblemMatchers struct {
	Matchers []GitHubMatcher `json:"problemMatcher,omitempty"`
}

// GitHubMatcher defines a problem matcher.
type GitHubMatcher struct {
	// Owner an ID field that can be used to remove or replace the problem matcher.
	// **required**
	Owner string `json:"owner,omitempty"`
	// Severity indicates the default severity, either 'warning' or 'error' case-insensitive.
	// Defaults to 'error'.
	Severity string          `json:"severity,omitempty"`
	Pattern  []GitHubPattern `json:"pattern,omitempty"`
}

// GitHubPattern defines a pattern for a problem matcher.
type GitHubPattern struct {
	// Regexp the regexp pattern that provides the groups to match against.
	// **required**
	Regexp string `json:"regexp,omitempty"`
	// File a group number containing the file name.
	File int `json:"file,omitempty"`
	// FromPath a group number containing a filepath used to root the file (e.g. a project file).
	FromPath int `json:"fromPath,omitempty"`
	// Line a group number containing the line number.
	Line int `json:"line,omitempty"`
	// Column a group number containing the column information.
	Column int `json:"column,omitempty"`
	// Severity a group number containing either 'warning' or 'error' case-insensitive.
	// Defaults to `error`.
	Severity int `json:"severity,omitempty"`
	// Code a group number containing the error code.
	Code int `json:"code,omitempty"`
	// Message a group number containing the error message.
	// **required** at least one pattern must set the message.
	Message int `json:"message,omitempty"`
	// Loop whether to loop until a match is not found,
	// only valid on the last pattern of a multi-pattern matcher.
	Loop bool `json:"loop,omitempty"`
}

type GitHub struct {
	tempPath string
	w        io.Writer
}

// NewGitHub output format outputs issues according to GitHub actions the problem matcher regexp.
func NewGitHub(w io.Writer) *GitHub {
	return &GitHub{
		tempPath: filepath.Join(os.TempDir(), filenameGitHubActionProblemMatchers),
		w:        w,
	}
}

func (p *GitHub) Print(issues []result.Issue) error {
	// Note: the file with the problem matcher definition should not be removed.
	// A sleep can mitigate this problem but this will be flaky.
	//
	// Result if the file is removed prematurely:
	// Error: Unable to process command '::add-matcher::/tmp/golangci-lint-action-problem-matchers.json' successfully.
	// Error: Could not find file '/tmp/golangci-lint-action-problem-matchers.json'.
	filename, err := p.storeProblemMatcher()
	if err != nil {
		return err
	}

	_, _ = fmt.Fprintln(p.w, "::debug::problem matcher definition file: "+filename)

	_, _ = fmt.Fprintln(p.w, "::add-matcher::"+filename)

	for ind := range issues {
		_, err := fmt.Fprintln(p.w, formatIssueAsGitHub(&issues[ind]))
		if err != nil {
			return err
		}
	}

	_, _ = fmt.Fprintln(p.w, "::remove-matcher owner=golangci-lint-action::")

	return nil
}

func (p *GitHub) storeProblemMatcher() (string, error) {
	file, err := os.Create(p.tempPath)
	if err != nil {
		return "", err
	}

	defer file.Close()

	err = json.NewEncoder(file).Encode(generateProblemMatcher())
	if err != nil {
		return "", err
	}

	return file.Name(), nil
}

func generateProblemMatcher() GitHubProblemMatchers {
	return GitHubProblemMatchers{
		Matchers: []GitHubMatcher{
			{
				Owner:    "golangci-lint-action",
				Severity: "error",
				Pattern: []GitHubPattern{
					{
						Regexp:   `^([^\s]+)\s+([^:]+):(\d+):(?:(\d+):)?\s+(.+)$`,
						Severity: 1,
						File:     2,
						Line:     3,
						Column:   4,
						Message:  5,
					},
				},
			},
		},
	}
}

func formatIssueAsGitHub(issue *result.Issue) string {
	severity := defaultGitHubSeverity
	if issue.Severity != "" {
		severity = issue.Severity
	}

	// Convert backslashes to forward slashes.
	// This is needed when running on windows.
	// Otherwise, GitHub won't be able to show the annotations pointing to the file path with backslashes.
	file := filepath.ToSlash(issue.FilePath())

	ret := fmt.Sprintf("%s\t%s:%d:", severity, file, issue.Line())
	if issue.Pos.Column != 0 {
		ret += fmt.Sprintf("%d:", issue.Pos.Column)
	}

	ret += fmt.Sprintf("\t%s (%s)", issue.Text, issue.FromLinter)
	return ret
}
