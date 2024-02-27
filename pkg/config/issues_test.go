package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetExcludePatterns(t *testing.T) {
	patterns := GetExcludePatterns(nil)

	assert.Equal(t, DefaultExcludePatterns, patterns)
}

func TestGetExcludePatterns_includes(t *testing.T) {
	include := []string{DefaultExcludePatterns[0].ID, DefaultExcludePatterns[1].ID}

	exclude := GetExcludePatterns(include)
	assert.Len(t, exclude, len(DefaultExcludePatterns)-len(include))

	for _, p := range exclude {
		assert.NotContains(t, include, p.ID)
		assert.Contains(t, DefaultExcludePatterns, p)
	}
}

func TestExcludeRule_Validate(t *testing.T) {
	testCases := []struct {
		desc     string
		rule     *ExcludeRule
		expected string
	}{
		{
			desc:     "empty rule",
			rule:     &ExcludeRule{},
			expected: "at least 2 of (text, source, path[-except],  linters) should be set",
		},
		{
			desc: "only path rule",
			rule: &ExcludeRule{
				BaseRule{
					Path: "test",
				},
			},
			expected: "at least 2 of (text, source, path[-except],  linters) should be set",
		},
		{
			desc: "only path-except rule",
			rule: &ExcludeRule{
				BaseRule{
					PathExcept: "test",
				},
			},
			expected: "at least 2 of (text, source, path[-except],  linters) should be set",
		},
		{
			desc: "only text rule",
			rule: &ExcludeRule{
				BaseRule{
					Text: "test",
				},
			},
			expected: "at least 2 of (text, source, path[-except],  linters) should be set",
		},
		{
			desc: "only source rule",
			rule: &ExcludeRule{
				BaseRule{
					Source: "test",
				},
			},
			expected: "at least 2 of (text, source, path[-except],  linters) should be set",
		},
		{
			desc: "invalid path rule",
			rule: &ExcludeRule{
				BaseRule{
					Path: "**test",
				},
			},
			expected: "invalid path regex: error parsing regexp: missing argument to repetition operator: `*`",
		},
		{
			desc: "invalid path-except rule",
			rule: &ExcludeRule{
				BaseRule{
					PathExcept: "**test",
				},
			},
			expected: "invalid path-except regex: error parsing regexp: missing argument to repetition operator: `*`",
		},
		{
			desc: "invalid text rule",
			rule: &ExcludeRule{
				BaseRule{
					Text: "**test",
				},
			},
			expected: "invalid text regex: error parsing regexp: missing argument to repetition operator: `*`",
		},
		{
			desc: "invalid source rule",
			rule: &ExcludeRule{
				BaseRule{
					Source: "**test",
				},
			},
			expected: "invalid source regex: error parsing regexp: missing argument to repetition operator: `*`",
		},
		{
			desc: "path and path-expect",
			rule: &ExcludeRule{
				BaseRule{
					Path:       "test",
					PathExcept: "test",
				},
			},
			expected: "path and path-except should not be set at the same time",
		},
	}

	for _, test := range testCases {
		test := test
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()

			err := test.rule.Validate()
			require.EqualError(t, err, test.expected)
		})
	}
}

func TestExcludeRule_Validate_error(t *testing.T) {
	testCases := []struct {
		desc string
		rule *ExcludeRule
	}{
		{
			desc: "path and linter",
			rule: &ExcludeRule{
				BaseRule{
					Path:    "test",
					Linters: []string{"a"},
				},
			},
		},
		{
			desc: "path-except and linter",
			rule: &ExcludeRule{
				BaseRule{
					PathExcept: "test",
					Linters:    []string{"a"},
				},
			},
		},
		{
			desc: "text and linter",
			rule: &ExcludeRule{
				BaseRule{
					Text:    "test",
					Linters: []string{"a"},
				},
			},
		},
		{
			desc: "source and linter",
			rule: &ExcludeRule{
				BaseRule{
					Source:  "test",
					Linters: []string{"a"},
				},
			},
		},
		{
			desc: "path and text",
			rule: &ExcludeRule{
				BaseRule{
					Path: "test",
					Text: "test",
				},
			},
		},
		{
			desc: "path and text and linter",
			rule: &ExcludeRule{
				BaseRule{
					Path:    "test",
					Text:    "test",
					Linters: []string{"a"},
				},
			},
		},
	}

	for _, test := range testCases {
		test := test
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()

			err := test.rule.Validate()
			require.NoError(t, err)
		})
	}
}
