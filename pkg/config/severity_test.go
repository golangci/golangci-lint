package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSeverity_Validate(t *testing.T) {
	rule := &SeverityRule{
		BaseRule: BaseRule{
			Path: "test",
		},
	}

	err := rule.Validate()
	require.NoError(t, err)
}

func TestSeverity_Validate_error(t *testing.T) {
	testCases := []struct {
		desc     string
		rule     *SeverityRule
		expected string
	}{
		{
			desc:     "empty rule",
			rule:     &SeverityRule{},
			expected: "at least 1 of (text, source, path[-except],  linters) should be set",
		},
		{
			desc: "invalid path rule",
			rule: &SeverityRule{
				BaseRule: BaseRule{
					Path: "**test",
				},
			},
			expected: "invalid path regex: error parsing regexp: missing argument to repetition operator: `*`",
		},
		{
			desc: "invalid path-except rule",
			rule: &SeverityRule{
				BaseRule: BaseRule{
					PathExcept: "**test",
				},
			},
			expected: "invalid path-except regex: error parsing regexp: missing argument to repetition operator: `*`",
		},
		{
			desc: "invalid text rule",
			rule: &SeverityRule{
				BaseRule: BaseRule{
					Text: "**test",
				},
			},
			expected: "invalid text regex: error parsing regexp: missing argument to repetition operator: `*`",
		},
		{
			desc: "invalid source rule",
			rule: &SeverityRule{
				BaseRule: BaseRule{
					Source: "**test",
				},
			},
			expected: "invalid source regex: error parsing regexp: missing argument to repetition operator: `*`",
		},
		{
			desc: "path and path-expect",
			rule: &SeverityRule{
				BaseRule: BaseRule{
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
