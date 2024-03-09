package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSeverity_Validate(t *testing.T) {
	testCases := []struct {
		desc     string
		severity *Severity
	}{
		{
			desc: "default with rules",
			severity: &Severity{
				Default: "high",
				Rules: []SeverityRule{
					{
						Severity: "low",
						BaseRule: BaseRule{
							Path: "test",
						},
					},
				},
			},
		},
		{
			desc: "default without rules",
			severity: &Severity{
				Default: "high",
			},
		},
		{
			desc: "same severity between default and rule",
			severity: &Severity{
				Default: "high",
				Rules: []SeverityRule{
					{
						Severity: "high",
						BaseRule: BaseRule{
							Path: "test",
						},
					},
				},
			},
		},
	}

	for _, test := range testCases {
		test := test
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()

			err := test.severity.Validate()
			require.NoError(t, err)
		})
	}
}

func TestSeverity_Validate_error(t *testing.T) {
	testCases := []struct {
		desc     string
		severity *Severity
		expected string
	}{
		{
			desc: "missing default severity",
			severity: &Severity{
				Default: "",
				Rules: []SeverityRule{
					{
						Severity: "low",
						BaseRule: BaseRule{
							Path: "test",
						},
					},
				},
			},
			expected: "can't set severity rule option: no default severity defined",
		},
		{
			desc: "missing rule severity",
			severity: &Severity{
				Default: "high",
				Rules: []SeverityRule{
					{
						BaseRule: BaseRule{
							Path: "test",
						},
					},
				},
			},
			expected: "error in severity rule #0: severity should be set",
		},
	}

	for _, test := range testCases {
		test := test
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()

			err := test.severity.Validate()
			require.EqualError(t, err, test.expected)
		})
	}
}

func TestSeverityRule_Validate(t *testing.T) {
	rule := &SeverityRule{
		Severity: "low",
		BaseRule: BaseRule{
			Path: "test",
		},
	}

	err := rule.Validate()
	require.NoError(t, err)
}

func TestSeverityRule_Validate_error(t *testing.T) {
	testCases := []struct {
		desc     string
		rule     *SeverityRule
		expected string
	}{
		{
			desc: "missing severity",
			rule: &SeverityRule{
				BaseRule: BaseRule{
					Path: "test",
				},
			},
			expected: "severity should be set",
		},
		{
			desc: "empty rule",
			rule: &SeverityRule{
				Severity: "low",
			},
			expected: "at least 1 of (text, source, path[-except],  linters) should be set",
		},
		{
			desc: "invalid path rule",
			rule: &SeverityRule{
				Severity: "low",
				BaseRule: BaseRule{
					Path: "**test",
				},
			},
			expected: "invalid path regex: error parsing regexp: missing argument to repetition operator: `*`",
		},
		{
			desc: "invalid path-except rule",
			rule: &SeverityRule{
				Severity: "low",
				BaseRule: BaseRule{
					PathExcept: "**test",
				},
			},
			expected: "invalid path-except regex: error parsing regexp: missing argument to repetition operator: `*`",
		},
		{
			desc: "invalid text rule",
			rule: &SeverityRule{
				Severity: "low",
				BaseRule: BaseRule{
					Text: "**test",
				},
			},
			expected: "invalid text regex: error parsing regexp: missing argument to repetition operator: `*`",
		},
		{
			desc: "invalid source rule",
			rule: &SeverityRule{
				Severity: "low",
				BaseRule: BaseRule{
					Source: "**test",
				},
			},
			expected: "invalid source regex: error parsing regexp: missing argument to repetition operator: `*`",
		},
		{
			desc: "path and path-expect",
			rule: &SeverityRule{
				Severity: "low",
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
