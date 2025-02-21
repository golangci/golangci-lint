package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLinterExclusions_Validate(t *testing.T) {
	testCases := []struct {
		desc       string
		exclusions *LinterExclusions
	}{
		{
			desc:       "empty configuration",
			exclusions: &LinterExclusions{},
		},
		{
			desc: "valid preset",
			exclusions: &LinterExclusions{
				Presets: []string{ExclusionPresetComments},
			},
		},
	}

	for _, test := range testCases {
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()

			err := test.exclusions.Validate()
			require.NoError(t, err)
		})
	}
}

func TestLinterExclusions_Validate_error(t *testing.T) {
	testCases := []struct {
		desc       string
		exclusions *LinterExclusions
		expected   string
	}{
		{
			desc: "invalid preset name",
			exclusions: &LinterExclusions{
				Presets: []string{"foo"},
			},
			expected: "invalid preset: foo",
		},
		{
			desc: "invalid rule: empty rule",
			exclusions: &LinterExclusions{
				Rules: []ExcludeRule{{BaseRule: BaseRule{}}},
			},
			expected: "error in exclude rule #0: at least 2 of (text, source, path[-except], linters) should be set",
		},
	}

	for _, test := range testCases {
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()

			err := test.exclusions.Validate()
			require.EqualError(t, err, test.expected)
		})
	}
}
