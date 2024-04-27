package processors

import (
	"path"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/fsutils"
	"github.com/golangci/golangci-lint/pkg/logutils"
	"github.com/golangci/golangci-lint/pkg/result"
)

func TestSeverity_multiple(t *testing.T) {
	lineCache := fsutils.NewLineCache(fsutils.NewFileCache())
	files := fsutils.NewFiles(lineCache, "")
	log := logutils.NewStderrLog(logutils.DebugKeyEmpty)

	opts := &config.Severity{
		Default: "error",
		Rules: []config.SeverityRule{
			{
				Severity: "info",
				BaseRule: config.BaseRule{
					Text:    "^ssl$",
					Linters: []string{"gosec"},
				},
			},
			{
				Severity: "info",
				BaseRule: config.BaseRule{
					Linters: []string{"linter"},
					Path:    "e.go",
				},
			},
			{
				Severity: "info",
				BaseRule: config.BaseRule{
					Text: "^testonly$",
					Path: `_test\.go`,
				},
			},
			{
				Severity: "info",
				BaseRule: config.BaseRule{
					Text:       "^nontestonly$",
					PathExcept: `_test\.go`,
				},
			},
			{
				BaseRule: config.BaseRule{
					Source:  "^//go:generate ",
					Linters: []string{"lll"},
				},
			},
			{
				Severity: "info",
				BaseRule: config.BaseRule{
					Source: "^//go:dosomething",
				},
			},
			{
				Severity: "info",
				BaseRule: config.BaseRule{
					Linters: []string{"someotherlinter"},
				},
			},
			{
				Severity: "info",
				BaseRule: config.BaseRule{
					Linters: []string{"somelinter"},
				},
			},
			{
				Severity: "info",
			},
		},
	}

	p := NewSeverity(log, files, opts)

	cases := []issueTestCase{
		{Path: "ssl.go", Text: "ssl", Linter: "gosec"},
		{Path: "e.go", Text: "some", Linter: "linter"},
		{Path: "e_test.go", Text: "testonly", Linter: "testlinter"},
		{Path: "e.go", Text: "nontestonly", Linter: "testlinter"},
		{Path: "e_test.go", Text: "nontestonly", Linter: "testlinter"},
		{Path: filepath.Join("testdata", "exclude_rules.go"), Line: 3, Linter: "lll"},
		{Path: filepath.Join("testdata", "severity_rules.go"), Line: 3, Linter: "invalidgo"},
		{Path: "someotherlinter.go", Text: "someotherlinter", Linter: "someotherlinter"},
		{Path: "somenotmatchlinter.go", Text: "somenotmatchlinter", Linter: "somenotmatchlinter"},
		{Path: "empty.go", Text: "empty", Linter: "empty"},
	}

	var issues []result.Issue
	for _, c := range cases {
		issues = append(issues, newIssueFromIssueTestCase(c))
	}

	processedIssues := process(t, p, issues...)

	var resultingCases []issueTestCase
	for _, i := range processedIssues {
		resultingCases = append(resultingCases, issueTestCase{
			Path:     i.FilePath(),
			Linter:   i.FromLinter,
			Text:     i.Text,
			Line:     i.Line(),
			Severity: i.Severity,
		})
	}

	expectedCases := []issueTestCase{
		{Path: "ssl.go", Text: "ssl", Linter: "gosec", Severity: "info"},
		{Path: "e.go", Text: "some", Linter: "linter", Severity: "info"},
		{Path: "e_test.go", Text: "testonly", Linter: "testlinter", Severity: "info"},
		{Path: "e.go", Text: "nontestonly", Linter: "testlinter", Severity: "info"},       // matched
		{Path: "e_test.go", Text: "nontestonly", Linter: "testlinter", Severity: "error"}, // not matched
		{Path: filepath.Join("testdata", "exclude_rules.go"), Line: 3, Linter: "lll", Severity: "error"},
		{Path: filepath.Join("testdata", "severity_rules.go"), Line: 3, Linter: "invalidgo", Severity: "info"},
		{Path: "someotherlinter.go", Text: "someotherlinter", Linter: "someotherlinter", Severity: "info"},
		{Path: "somenotmatchlinter.go", Text: "somenotmatchlinter", Linter: "somenotmatchlinter", Severity: "error"},
		{Path: "empty.go", Text: "empty", Linter: "empty", Severity: "error"},
	}

	assert.Equal(t, expectedCases, resultingCases)
}

func TestSeverity_pathPrefix(t *testing.T) {
	lineCache := fsutils.NewLineCache(fsutils.NewFileCache())
	pathPrefix := path.Join("some", "dir")
	files := fsutils.NewFiles(lineCache, pathPrefix)
	log := logutils.NewStderrLog(logutils.DebugKeyEmpty)

	opts := &config.Severity{
		Default: "error",
		Rules: []config.SeverityRule{
			{
				Severity: "info",
				BaseRule: config.BaseRule{
					Text: "some",
					Path: `some/dir/e\.go`,
				},
			},
		},
	}

	p := NewSeverity(log, files, opts)

	cases := []issueTestCase{
		{Path: "e.go", Text: "some", Linter: "linter"},
		{Path: "other.go", Text: "some", Linter: "linter"},
	}

	var issues []result.Issue
	for _, c := range cases {
		issues = append(issues, newIssueFromIssueTestCase(c))
	}

	processedIssues := process(t, p, issues...)

	var resultingCases []issueTestCase
	for _, i := range processedIssues {
		resultingCases = append(resultingCases, issueTestCase{
			Path:     i.FilePath(),
			Linter:   i.FromLinter,
			Text:     i.Text,
			Line:     i.Line(),
			Severity: i.Severity,
		})
	}

	expectedCases := []issueTestCase{
		{Path: "e.go", Text: "some", Linter: "linter", Severity: "info"},
		{Path: "other.go", Text: "some", Linter: "linter", Severity: "error"},
	}

	assert.Equal(t, expectedCases, resultingCases)
}

func TestSeverity_text(t *testing.T) {
	opts := &config.Severity{
		Rules: []config.SeverityRule{
			{
				BaseRule: config.BaseRule{
					Text:    "^severity$",
					Linters: []string{"linter"},
				},
			},
		},
	}

	p := NewSeverity(nil, nil, opts)

	texts := []string{"seveRity", "1", "", "serverit", "notseverity"}
	var issues []result.Issue
	for _, t := range texts {
		issues = append(issues, result.Issue{
			Text:       t,
			FromLinter: "linter",
		})
	}

	processedIssues := process(t, p, issues...)
	assert.Len(t, processedIssues, len(issues))

	var processedTexts []string
	for _, i := range processedIssues {
		processedTexts = append(processedTexts, i.Text)
	}

	assert.Equal(t, texts, processedTexts)
}

func TestSeverity_onlyDefault(t *testing.T) {
	lineCache := fsutils.NewLineCache(fsutils.NewFileCache())
	files := fsutils.NewFiles(lineCache, "")
	log := logutils.NewStderrLog(logutils.DebugKeyEmpty)

	opts := config.Severity{
		Default: "info",
		Rules:   []config.SeverityRule{},
	}

	p := NewSeverity(log, files, &opts)

	cases := []issueTestCase{
		{Path: "ssl.go", Text: "ssl", Linter: "gosec"},
		{Path: "empty.go", Text: "empty", Linter: "empty"},
	}

	var issues []result.Issue
	for _, c := range cases {
		issues = append(issues, newIssueFromIssueTestCase(c))
	}

	processedIssues := process(t, p, issues...)

	var resultingCases []issueTestCase
	for _, i := range processedIssues {
		resultingCases = append(resultingCases, issueTestCase{
			Path:     i.FilePath(),
			Linter:   i.FromLinter,
			Text:     i.Text,
			Line:     i.Line(),
			Severity: i.Severity,
		})
	}

	expectedCases := []issueTestCase{
		{Path: "ssl.go", Text: "ssl", Linter: "gosec", Severity: "info"},
		{Path: "empty.go", Text: "empty", Linter: "empty", Severity: "info"},
	}

	assert.Equal(t, expectedCases, resultingCases)
}

func TestSeverity_empty(t *testing.T) {
	p := NewSeverity(nil, nil, &config.Severity{})

	processAssertSame(t, p, newIssueFromTextTestCase("test"))
}

func TestSeverity_caseSensitive(t *testing.T) {
	lineCache := fsutils.NewLineCache(fsutils.NewFileCache())
	files := fsutils.NewFiles(lineCache, "")

	opts := &config.Severity{
		Default: "error",
		Rules: []config.SeverityRule{
			{
				Severity: "info",
				BaseRule: config.BaseRule{
					Text:    "^ssl$",
					Linters: []string{"gosec", "someotherlinter"},
				},
			},
		},
		CaseSensitive: true,
	}

	p := NewSeverity(nil, files, opts)

	cases := []issueTestCase{
		{Path: "e.go", Text: "ssL", Linter: "gosec"},
	}

	var issues []result.Issue
	for _, c := range cases {
		issues = append(issues, newIssueFromIssueTestCase(c))
	}

	processedIssues := process(t, p, issues...)

	var resultingCases []issueTestCase
	for _, i := range processedIssues {
		resultingCases = append(resultingCases, issueTestCase{
			Path:     i.FilePath(),
			Linter:   i.FromLinter,
			Text:     i.Text,
			Line:     i.Line(),
			Severity: i.Severity,
		})
	}

	expectedCases := []issueTestCase{
		{Path: "e.go", Text: "ssL", Linter: "gosec", Severity: "error"},
	}

	assert.Equal(t, expectedCases, resultingCases)
}

func TestSeverity_transform(t *testing.T) {
	lineCache := fsutils.NewLineCache(fsutils.NewFileCache())
	files := fsutils.NewFiles(lineCache, "")

	testCases := []struct {
		desc     string
		opts     *config.Severity
		issue    *result.Issue
		expected *result.Issue
	}{
		{
			desc: "apply severity from rule",
			opts: &config.Severity{
				Default: "error",
				Rules: []config.SeverityRule{
					{
						Severity: "info",
						BaseRule: config.BaseRule{
							Linters: []string{"linter1"},
						},
					},
				},
			},
			issue: &result.Issue{
				Text:       "This is a report",
				FromLinter: "linter1",
			},
			expected: &result.Issue{
				Text:       "This is a report",
				FromLinter: "linter1",
				Severity:   "info",
			},
		},
		{
			desc: "apply severity from default",
			opts: &config.Severity{
				Default: "error",
				Rules: []config.SeverityRule{
					{
						Severity: "info",
						BaseRule: config.BaseRule{
							Linters: []string{"linter1"},
						},
					},
				},
			},
			issue: &result.Issue{
				Text:       "This is a report",
				FromLinter: "linter2",
			},
			expected: &result.Issue{
				Text:       "This is a report",
				FromLinter: "linter2",
				Severity:   "error",
			},
		},
		{
			desc: "severity from rule override severity from linter",
			opts: &config.Severity{
				Default: "error",
				Rules: []config.SeverityRule{
					{
						Severity: "info",
						BaseRule: config.BaseRule{
							Linters: []string{"linter1"},
						},
					},
				},
			},
			issue: &result.Issue{
				Text:       "This is a report",
				FromLinter: "linter1",
				Severity:   "huge",
			},
			expected: &result.Issue{
				Text:       "This is a report",
				FromLinter: "linter1",
				Severity:   "info",
			},
		},
		{
			desc: "severity from default override severity from linter",
			opts: &config.Severity{
				Default: "error",
				Rules: []config.SeverityRule{
					{
						Severity: "info",
						BaseRule: config.BaseRule{
							Linters: []string{"linter1"},
						},
					},
				},
			},
			issue: &result.Issue{
				Text:       "This is a report",
				FromLinter: "linter2",
				Severity:   "huge",
			},
			expected: &result.Issue{
				Text:       "This is a report",
				FromLinter: "linter2",
				Severity:   "error",
			},
		},
		{
			desc: "keep severity from linter as rule",
			opts: &config.Severity{
				Default: "error",
				Rules: []config.SeverityRule{
					{
						Severity: severityFromLinter,
						BaseRule: config.BaseRule{
							Linters: []string{"linter1"},
						},
					},
				},
			},
			issue: &result.Issue{
				Text:       "This is a report",
				FromLinter: "linter1",
				Severity:   "huge",
			},
			expected: &result.Issue{
				Text:       "This is a report",
				FromLinter: "linter1",
				Severity:   "huge",
			},
		},
		{
			desc: "keep severity from linter as default",
			opts: &config.Severity{
				Default: severityFromLinter,
				Rules: []config.SeverityRule{
					{
						Severity: "info",
						BaseRule: config.BaseRule{
							Linters: []string{"linter1"},
						},
					},
				},
			},
			issue: &result.Issue{
				Text:       "This is a report",
				FromLinter: "linter2",
				Severity:   "huge",
			},
			expected: &result.Issue{
				Text:       "This is a report",
				FromLinter: "linter2",
				Severity:   "huge",
			},
		},
		{
			desc: "keep severity from linter as default (without rule)",
			opts: &config.Severity{
				Default: severityFromLinter,
			},
			issue: &result.Issue{
				Text:       "This is a report",
				FromLinter: "linter2",
				Severity:   "huge",
			},
			expected: &result.Issue{
				Text:       "This is a report",
				FromLinter: "linter2",
				Severity:   "huge",
			},
		},
	}

	for _, test := range testCases {
		test := test
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()

			p := NewSeverity(nil, files, test.opts)

			newIssue := p.transform(test.issue)

			assert.Equal(t, test.expected, newIssue)
		})
	}
}
