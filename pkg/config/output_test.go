package config

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/golangci/golangci-lint/v2/pkg/fsutils"
)

func TestOutput_Validate(t *testing.T) {
	testCases := []struct {
		desc     string
		settings *Output
	}{
		{
			desc: "SortOrder: file",
			settings: &Output{
				SortOrder: []string{"file"},
			},
		},
		{
			desc: "SortOrder: linter",
			settings: &Output{
				SortOrder: []string{"linter"},
			},
		},
		{
			desc: "SortOrder: severity",
			settings: &Output{
				SortOrder: []string{"severity"},
			},
		},
		{
			desc: "SortOrder: multiple",
			settings: &Output{
				SortOrder: []string{"file", "linter", "severity"},
			},
		},
		{
			desc: "PathMode: empty",
			settings: &Output{
				PathMode: "",
			},
		},
		{
			desc: "PathMode: absolute",
			settings: &Output{
				PathMode: fsutils.OutputPathModeAbsolute,
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
			desc: "SortOrder: invalid",
			settings: &Output{
				SortOrder: []string{"a"},
			},
			expected: `unsupported sort-order name "a"`,
		},
		{
			desc: "SortOrder: duplicate",
			settings: &Output{
				SortOrder: []string{"file", "linter", "severity", "linter"},
			},
			expected: `the sort-order name "linter" is repeated several times`,
		},
		{
			desc: "PathMode: invalid",
			settings: &Output{
				PathMode: "example",
			},
			expected: `unsupported output path mode "example"`,
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
