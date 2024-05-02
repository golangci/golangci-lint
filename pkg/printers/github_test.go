package printers

import (
	"bytes"
	"fmt"
	"go/token"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/golangci/golangci-lint/pkg/result"
)

func TestGitHub_Print(t *testing.T) {
	issues := []result.Issue{
		{
			FromLinter: "linter-a",
			Severity:   "warning",
			Text:       "some issue",
			Pos: token.Position{
				Filename: "path/to/filea.go",
				Offset:   2,
				Line:     10,
				Column:   4,
			},
		},
		{
			FromLinter: "linter-b",
			Severity:   "error",
			Text:       "another issue",
			SourceLines: []string{
				"func foo() {",
				"\tfmt.Println(\"bar\")",
				"}",
			},
			Pos: token.Position{
				Filename: "path/to/fileb.go",
				Offset:   5,
				Line:     300,
				Column:   9,
			},
		},
	}

	buf := new(bytes.Buffer)

	printer := NewGitHub(buf)
	printer.tempPath = filepath.Join(t.TempDir(), filenameGitHubActionProblemMatchers)

	err := printer.Print(issues)
	require.NoError(t, err)

	expected := `::debug::problem matcher definition file: /tmp/golangci-lint-action-problem-matchers.json
::add-matcher::/tmp/golangci-lint-action-problem-matchers.json
warning	path/to/filea.go:10:4:	some issue (linter-a)
error	path/to/fileb.go:300:9:	another issue (linter-b)
::remove-matcher owner=golangci-lint-action::
`
	// To support all the OS.
	expected = strings.ReplaceAll(expected, "/tmp/golangci-lint-action-problem-matchers.json", printer.tempPath)

	assert.Equal(t, expected, buf.String())
}

func Test_formatIssueAsGitHub(t *testing.T) {
	sampleIssue := result.Issue{
		FromLinter: "sample-linter",
		Text:       "some issue",
		Pos: token.Position{
			Filename: "path/to/file.go",
			Offset:   2,
			Line:     10,
			Column:   4,
		},
	}
	require.Equal(t, "error\tpath/to/file.go:10:4:\tsome issue (sample-linter)", formatIssueAsGitHub(&sampleIssue))

	sampleIssue.Pos.Column = 0
	require.Equal(t, "error\tpath/to/file.go:10:\tsome issue (sample-linter)", formatIssueAsGitHub(&sampleIssue))
}

func Test_formatIssueAsGitHub_Windows(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("Skipping test on non Windows")
	}

	sampleIssue := result.Issue{
		FromLinter: "sample-linter",
		Text:       "some issue",
		Pos: token.Position{
			Filename: "path\\to\\file.go",
			Offset:   2,
			Line:     10,
			Column:   4,
		},
	}
	require.Equal(t, "error\tpath/to/file.go:10:4:\tsome issue (sample-linter)", formatIssueAsGitHub(&sampleIssue))

	sampleIssue.Pos.Column = 0
	require.Equal(t, "error\tpath/to/file.go:10:\tsome issue (sample-linter)", formatIssueAsGitHub(&sampleIssue))
}

func Test_generateProblemMatcher(t *testing.T) {
	pattern := generateProblemMatcher().Matchers[0].Pattern[0]

	exp := regexp.MustCompile(pattern.Regexp)

	testCases := []struct {
		desc     string
		line     string
		expected string
	}{
		{
			desc: "error",
			line: "error\tpath/to/filea.go:10:4:\tsome issue (sample-linter)",
			expected: `File: path/to/filea.go
Line: 10
Column: 4
Severity: error
Message: some issue (sample-linter)`,
		},
		{
			desc: "warning",
			line: "warning\tpath/to/fileb.go:1:4:\tsome issue (sample-linter)",
			expected: `File: path/to/fileb.go
Line: 1
Column: 4
Severity: warning
Message: some issue (sample-linter)`,
		},
		{
			desc: "no column",
			line: "error\t \tpath/to/fileb.go:40:\t Foo bar",
			expected: `File: path/to/fileb.go
Line: 40
Column: 
Severity: error
Message: Foo bar`,
		},
	}

	for _, test := range testCases {
		test := test
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()

			assert.True(t, exp.MatchString(test.line), test.line)

			actual := exp.ReplaceAllString(test.line, createReplacement(&pattern))

			assert.Equal(t, test.expected, actual)
		})
	}
}

func createReplacement(pattern *GitHubPattern) string {
	var repl []string

	if pattern.File > 0 {
		repl = append(repl, fmt.Sprintf("File: $%d", pattern.File))
	}

	if pattern.FromPath > 0 {
		repl = append(repl, fmt.Sprintf("FromPath: $%d", pattern.FromPath))
	}

	if pattern.Line > 0 {
		repl = append(repl, fmt.Sprintf("Line: $%d", pattern.Line))
	}

	if pattern.Column > 0 {
		repl = append(repl, fmt.Sprintf("Column: $%d", pattern.Column))
	}

	if pattern.Severity > 0 {
		repl = append(repl, fmt.Sprintf("Severity: $%d", pattern.Severity))
	}

	if pattern.Code > 0 {
		repl = append(repl, fmt.Sprintf("Code: $%d", pattern.Code))
	}

	if pattern.Message > 0 {
		repl = append(repl, fmt.Sprintf("Message: $%d", pattern.Message))
	}

	if pattern.Loop {
		repl = append(repl, fmt.Sprintf("Loop: $%v", pattern.Loop))
	}

	return strings.Join(repl, "\n")
}
