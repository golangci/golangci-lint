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
				SortResults: true,
				SortOrder:   []string{"file"},
			},
		},
		{
			desc: "linter",
			settings: &Output{
				SortResults: true,
				SortOrder:   []string{"linter"},
			},
		},
		{
			desc: "severity",
			settings: &Output{
				SortResults: true,
				SortOrder:   []string{"severity"},
			},
		},
		{
			desc: "multiple",
			settings: &Output{
				SortResults: true,
				SortOrder:   []string{"file", "linter", "severity"},
			},
		},
	}

	for _, test := range testCases {
		test := test
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
			desc: "sort-results false and sort-order",
			settings: &Output{
				SortOrder: []string{"file"},
			},
			expected: "sort-results should be 'true' to use sort-order",
		},
		{
			desc: "invalid sort-order",
			settings: &Output{
				SortResults: true,
				SortOrder:   []string{"a"},
			},
			expected: `unsupported sort-order name "a"`,
		},
		{
			desc: "duplicate",
			settings: &Output{
				SortResults: true,
				SortOrder:   []string{"file", "linter", "severity", "linter"},
			},
			expected: `the sort-order name "linter" is repeated several times`,
		},
		{
			desc: "unsupported format",
			settings: &Output{
				Formats: []OutputFormat{
					{
						Format: "test",
					},
				},
			},
			expected: `unsupported output format "test"`,
		},
	}

	for _, test := range testCases {
		test := test
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()

			err := test.settings.Validate()
			require.EqualError(t, err, test.expected)
		})
	}
}

func TestOutputFormat_Validate(t *testing.T) {
	testCases := []struct {
		desc     string
		settings *OutputFormat
	}{
		{
			desc: "only format",
			settings: &OutputFormat{
				Format: "json",
			},
		},
		{
			desc: "format and path (relative)",
			settings: &OutputFormat{
				Format: "json",
				Path:   "./example.json",
			},
		},
		{
			desc: "format and path (absolute)",
			settings: &OutputFormat{
				Format: "json",
				Path:   "/tmp/example.json",
			},
		},
	}

	for _, test := range testCases {
		test := test
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()

			err := test.settings.Validate()
			require.NoError(t, err)
		})
	}
}

func TestOutputFormat_Validate_error(t *testing.T) {
	testCases := []struct {
		desc     string
		settings *OutputFormat
		expected string
	}{
		{
			desc:     "empty",
			settings: &OutputFormat{},
			expected: "the format is required",
		},
		{
			desc: "unsupported format",
			settings: &OutputFormat{
				Format: "test",
			},
			expected: `unsupported output format "test"`,
		},
	}

	for _, test := range testCases {
		test := test
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()

			err := test.settings.Validate()
			require.EqualError(t, err, test.expected)
		})
	}
}
