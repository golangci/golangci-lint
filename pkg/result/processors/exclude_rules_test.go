package processors

import (
	"path"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/fsutils"
	"github.com/golangci/golangci-lint/pkg/result"
)

func TestExcludeRules_multiple(t *testing.T) {
	lineCache := fsutils.NewLineCache(fsutils.NewFileCache())
	files := fsutils.NewFiles(lineCache, "")

	opts := &config.Issues{ExcludeRules: []config.ExcludeRule{
		{
			BaseRule: config.BaseRule{
				Text:    "^exclude$",
				Linters: []string{"linter"},
			},
		},
		{
			BaseRule: config.BaseRule{
				Linters: []string{"testlinter"},
				Path:    `_test\.go`,
			},
		},
		{
			BaseRule: config.BaseRule{
				Text: "^testonly$",
				Path: `_test\.go`,
			},
		},
		{
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
	}}

	p := NewExcludeRules(nil, files, opts)

	cases := []issueTestCase{
		{Path: "e.go", Text: "exclude", Linter: "linter"},
		{Path: "e.go", Text: "some", Linter: "linter"},
		{Path: "e_test.go", Text: "normal", Linter: "testlinter"},
		{Path: "e_Test.go", Text: "normal", Linter: "testlinter"},
		{Path: "e_test.go", Text: "another", Linter: "linter"},
		{Path: "e_test.go", Text: "testonly", Linter: "linter"},
		{Path: "e.go", Text: "nontestonly", Linter: "linter"},
		{Path: "e_test.go", Text: "nontestonly", Linter: "linter"},
		{Path: filepath.Join("testdata", "exclude_rules.go"), Line: 3, Linter: "lll"},
	}

	var issues []result.Issue
	for _, c := range cases {
		issues = append(issues, newIssueFromIssueTestCase(c))
	}

	processedIssues := process(t, p, issues...)

	var resultingCases []issueTestCase
	for _, i := range processedIssues {
		resultingCases = append(resultingCases, issueTestCase{
			Path:   i.FilePath(),
			Linter: i.FromLinter,
			Text:   i.Text,
			Line:   i.Line(),
		})
	}

	expectedCases := []issueTestCase{
		{Path: "e.go", Text: "some", Linter: "linter"},
		{Path: "e_Test.go", Text: "normal", Linter: "testlinter"},
		{Path: "e_test.go", Text: "another", Linter: "linter"},
		{Path: "e_test.go", Text: "nontestonly", Linter: "linter"},
	}

	assert.Equal(t, expectedCases, resultingCases)
}

func TestExcludeRules_pathPrefix(t *testing.T) {
	lineCache := fsutils.NewLineCache(fsutils.NewFileCache())
	pathPrefix := path.Join("some", "dir")
	files := fsutils.NewFiles(lineCache, pathPrefix)

	opts := &config.Issues{
		ExcludeRules: []config.ExcludeRule{
			{
				BaseRule: config.BaseRule{
					Path: `some/dir/e\.go`,
				},
			},
		},
	}

	p := NewExcludeRules(nil, files, opts)

	cases := []issueTestCase{
		{Path: "e.go"},
		{Path: "other.go"},
	}

	var issues []result.Issue
	for _, c := range cases {
		issues = append(issues, newIssueFromIssueTestCase(c))
	}

	processedIssues := process(t, p, issues...)

	var resultingCases []issueTestCase
	for _, i := range processedIssues {
		resultingCases = append(resultingCases, issueTestCase{
			Path:   i.FilePath(),
			Linter: i.FromLinter,
			Text:   i.Text,
			Line:   i.Line(),
		})
	}

	expectedCases := []issueTestCase{
		{Path: "other.go"},
	}

	assert.Equal(t, expectedCases, resultingCases)
}

func TestExcludeRules_text(t *testing.T) {
	opts := &config.Issues{
		ExcludeRules: []config.ExcludeRule{
			{
				BaseRule: config.BaseRule{
					Text:    "^exclude$",
					Linters: []string{"linter"},
				},
			},
		},
	}

	p := NewExcludeRules(nil, nil, opts)

	texts := []string{"excLude", "1", "", "exclud", "notexclude"}
	var issues []result.Issue
	for _, t := range texts {
		issues = append(issues, result.Issue{
			Text:       t,
			FromLinter: "linter",
		})
	}

	processedIssues := process(t, p, issues...)
	assert.Len(t, processedIssues, len(issues)-1)

	var processedTexts []string
	for _, i := range processedIssues {
		processedTexts = append(processedTexts, i.Text)
	}

	assert.Equal(t, texts[1:], processedTexts)
}

func TestExcludeRules_empty(t *testing.T) {
	processAssertSame(t, NewExcludeRules(nil, nil, &config.Issues{}), newIssueFromTextTestCase("test"))
}

func TestExcludeRules_caseSensitive_multiple(t *testing.T) {
	lineCache := fsutils.NewLineCache(fsutils.NewFileCache())
	files := fsutils.NewFiles(lineCache, "")

	opts := &config.Issues{
		ExcludeCaseSensitive: true,
		ExcludeRules: []config.ExcludeRule{
			{
				BaseRule: config.BaseRule{
					Text:    "^exclude$",
					Linters: []string{"linter"},
				},
			},
			{
				BaseRule: config.BaseRule{
					Linters: []string{"testlinter"},
					Path:    `_test\.go`,
				},
			},
			{
				BaseRule: config.BaseRule{
					Text: "^testonly$",
					Path: `_test\.go`,
				},
			},
			{
				BaseRule: config.BaseRule{
					Source:  "^//go:generate ",
					Linters: []string{"lll"},
				},
			},
		},
	}

	p := NewExcludeRules(nil, files, opts)

	cases := []issueTestCase{
		{Path: "e.go", Text: "exclude", Linter: "linter"},
		{Path: "e.go", Text: "excLude", Linter: "linter"},
		{Path: "e.go", Text: "some", Linter: "linter"},
		{Path: "e_test.go", Text: "normal", Linter: "testlinter"},
		{Path: "e_Test.go", Text: "normal", Linter: "testlinter"},
		{Path: "e_test.go", Text: "another", Linter: "linter"},
		{Path: "e_test.go", Text: "testonly", Linter: "linter"},
		{Path: "e_test.go", Text: "testOnly", Linter: "linter"},
		{Path: filepath.Join("testdata", "exclude_rules_case_sensitive.go"), Line: 3, Linter: "lll"},
	}

	var issues []result.Issue
	for _, c := range cases {
		issues = append(issues, newIssueFromIssueTestCase(c))
	}

	processedIssues := process(t, p, issues...)

	var resultingCases []issueTestCase
	for _, i := range processedIssues {
		resultingCases = append(resultingCases, issueTestCase{
			Path:   i.FilePath(),
			Linter: i.FromLinter,
			Text:   i.Text,
			Line:   i.Line(),
		})
	}

	expectedCases := []issueTestCase{
		{Path: "e.go", Text: "excLude", Linter: "linter"},
		{Path: "e.go", Text: "some", Linter: "linter"},
		{Path: "e_Test.go", Text: "normal", Linter: "testlinter"},
		{Path: "e_test.go", Text: "another", Linter: "linter"},
		{Path: "e_test.go", Text: "testOnly", Linter: "linter"},
		{Path: filepath.Join("testdata", "exclude_rules_case_sensitive.go"), Line: 3, Linter: "lll"},
	}

	assert.Equal(t, expectedCases, resultingCases)
}

func TestExcludeRules_caseSensitive_text(t *testing.T) {
	opts := &config.Issues{
		ExcludeCaseSensitive: true,
		ExcludeRules: []config.ExcludeRule{
			{
				BaseRule: config.BaseRule{
					Text:    "^exclude$",
					Linters: []string{"linter"},
				},
			},
		},
	}

	p := NewExcludeRules(nil, nil, opts)

	texts := []string{"exclude", "excLude", "1", "", "exclud", "notexclude"}

	var issues []result.Issue
	for _, t := range texts {
		issues = append(issues, result.Issue{
			Text:       t,
			FromLinter: "linter",
		})
	}

	processedIssues := process(t, p, issues...)
	assert.Len(t, processedIssues, len(issues)-1)

	var processedTexts []string
	for _, i := range processedIssues {
		processedTexts = append(processedTexts, i.Text)
	}

	assert.Equal(t, texts[1:], processedTexts)
}

func TestExcludeRules_caseSensitive_empty(t *testing.T) {
	processAssertSame(t, NewExcludeRules(nil, nil, &config.Issues{ExcludeCaseSensitive: true}), newIssueFromTextTestCase("test"))
}
