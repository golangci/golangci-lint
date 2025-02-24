package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestOutput_Validate(t *testing.T) {
	testCases := []struct {
		desc     string
		settings *Output
	}{
		{
			desc: "file",
			settings: &Output{
				SortOrder: []string{"file"},
			},
		},
		{
			desc: "linter",
			settings: &Output{
				SortOrder: []string{"linter"},
			},
		},
		{
			desc: "severity",
			settings: &Output{
				SortOrder: []string{"severity"},
			},
		},
		{
			desc: "multiple",
			settings: &Output{
				SortOrder: []string{"file", "linter", "severity"},
			},
		},
	}

	for _, test := range testCases {
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()

			err := test.settings.Validate()
			require.NoError(t, err)
		})
	}
}

func TestOutput_Validate_error(t *testing.T) {
	testCases := []struct {
		desc     string
		settings *Output
		expected string
	}{
		{
			desc: "invalid sort-order",
			settings: &Output{
				SortOrder: []string{"a"},
			},
			expected: `unsupported sort-order name "a"`,
		},
		{
			desc: "duplicate",
			settings: &Output{
				SortOrder: []string{"file", "linter", "severity", "linter"},
			},
			expected: `the sort-order name "linter" is repeated several times`,
		},
	}

	for _, test := range testCases {
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()

			err := test.settings.Validate()
			require.EqualError(t, err, test.expected)
		})
	}
}
