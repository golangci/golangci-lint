package processors

import (
	"path"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/golangci/golangci-lint/pkg/fsutils"
	"github.com/golangci/golangci-lint/pkg/logutils"
	"github.com/golangci/golangci-lint/pkg/report"
	"github.com/golangci/golangci-lint/pkg/result"
)

func TestSeverityRulesMultiple(t *testing.T) {
	lineCache := fsutils.NewLineCache(fsutils.NewFileCache())
	files := fsutils.NewFiles(lineCache, "")
	log := report.NewLogWrapper(logutils.NewStderrLog(logutils.DebugKeyEmpty), &report.Data{})
	p := NewSeverityRules("error", []SeverityRule{
		{
			Severity: "info",
			BaseRule: BaseRule{
				Text:    "^ssl$",
				Linters: []string{"gosec"},
			},
		},
		{
			Severity: "info",
			BaseRule: BaseRule{
				Linters: []string{"linter"},
				Path:    "e.go",
			},
		},
		{
			Severity: "info",
			BaseRule: BaseRule{
				Text: "^testonly$",
				Path: `_test\.go`,
			},
		},
		{
			BaseRule: BaseRule{
				Source:  "^//go:generate ",
				Linters: []string{"lll"},
			},
		},
		{
			Severity: "info",
			BaseRule: BaseRule{
				Source: "^//go:dosomething",
			},
		},
		{
			Severity: "info",
			BaseRule: BaseRule{
				Linters: []string{"someotherlinter"},
			},
		},
		{
			Severity: "info",
			BaseRule: BaseRule{
				Linters: []string{"somelinter"},
			},
		},
		{
			Severity: "info",
		},
	}, files, log)

	cases := []issueTestCase{
		{Path: "ssl.go", Text: "ssl", Linter: "gosec"},
		{Path: "e.go", Text: "some", Linter: "linter"},
		{Path: "e_test.go", Text: "testonly", Linter: "testlinter"},
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
		{Path: filepath.Join("testdata", "exclude_rules.go"), Line: 3, Linter: "lll", Severity: "error"},
		{Path: filepath.Join("testdata", "severity_rules.go"), Line: 3, Linter: "invalidgo", Severity: "info"},
		{Path: "someotherlinter.go", Text: "someotherlinter", Linter: "someotherlinter", Severity: "info"},
		{Path: "somenotmatchlinter.go", Text: "somenotmatchlinter", Linter: "somenotmatchlinter", Severity: "error"},
		{Path: "empty.go", Text: "empty", Linter: "empty", Severity: "error"},
	}
	assert.Equal(t, expectedCases, resultingCases)
}

func TestSeverityRulesPathPrefix(t *testing.T) {
	lineCache := fsutils.NewLineCache(fsutils.NewFileCache())
	pathPrefix := path.Join("some", "dir")
	files := fsutils.NewFiles(lineCache, pathPrefix)
	log := report.NewLogWrapper(logutils.NewStderrLog(logutils.DebugKeyEmpty), &report.Data{})
	p := NewSeverityRules("error", []SeverityRule{
		{
			Severity: "info",
			BaseRule: BaseRule{
				Text: "some",
				Path: `some/dir/e\.go`,
			},
		},
	}, files, log)

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

func TestSeverityRulesText(t *testing.T) {
	p := NewSeverityRules("", []SeverityRule{
		{
			BaseRule: BaseRule{
				Text:    "^severity$",
				Linters: []string{"linter"},
			},
		},
	}, nil, nil)
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

func TestSeverityRulesOnlyDefault(t *testing.T) {
	lineCache := fsutils.NewLineCache(fsutils.NewFileCache())
	files := fsutils.NewFiles(lineCache, "")
	log := report.NewLogWrapper(logutils.NewStderrLog(logutils.DebugKeyEmpty), &report.Data{})
	p := NewSeverityRules("info", []SeverityRule{}, files, log)

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

func TestSeverityRulesEmpty(t *testing.T) {
	processAssertSame(t, NewSeverityRules("", nil, nil, nil), newIssueFromTextTestCase("test"))
}

func TestSeverityRulesCaseSensitive(t *testing.T) {
	lineCache := fsutils.NewLineCache(fsutils.NewFileCache())
	files := fsutils.NewFiles(lineCache, "")
	p := NewSeverityRulesCaseSensitive("error", []SeverityRule{
		{
			Severity: "info",
			BaseRule: BaseRule{
				Text:    "^ssl$",
				Linters: []string{"gosec", "someotherlinter"},
			},
		},
	}, files, nil)

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
