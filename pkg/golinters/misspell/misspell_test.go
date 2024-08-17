package misspell

import (
	"testing"

	"github.com/golangci/misspell"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/golangci/golangci-lint/pkg/config"
)

func Test_appendExtraWords(t *testing.T) {
	extraWords := []config.MisspellExtraWords{
		{
			Typo:       "iff",
			Correction: "if",
		},
		{
			Typo:       "canCELation",
			Correction: "canceLLaTION",
		},
	}

	replacer := &misspell.Replacer{}

	err := appendExtraWords(replacer, extraWords)
	require.NoError(t, err)

	expected := []string{"iff", "if", "cancelation", "cancellation"}

	assert.Equal(t, expected, replacer.Replacements)
}

func Test_appendExtraWords_error(t *testing.T) {
	testCases := []struct {
		desc       string
		extraWords []config.MisspellExtraWords
		expected   string
	}{
		{
			desc: "empty fields",
			extraWords: []config.MisspellExtraWords{{
				Typo:       "",
				Correction: "",
			}},
			expected: `typo ("") and correction ("") fields should not be empty`,
		},
		{
			desc: "empty typo",
			extraWords: []config.MisspellExtraWords{{
				Typo:       "",
				Correction: "if",
			}},
			expected: `typo ("") and correction ("if") fields should not be empty`,
		},
		{
			desc: "empty correction",
			extraWords: []config.MisspellExtraWords{{
				Typo:       "iff",
				Correction: "",
			}},
			expected: `typo ("iff") and correction ("") fields should not be empty`,
		},
		{
			desc: "invalid characters in typo",
			extraWords: []config.MisspellExtraWords{{
				Typo:       "i'ff",
				Correction: "if",
			}},
			expected: `the word "i'ff" in the 'typo' field should only contain letters`,
		},
		{
			desc: "invalid characters in correction",
			extraWords: []config.MisspellExtraWords{{
				Typo:       "iff",
				Correction: "i'f",
			}},
			expected: `the word "i'f" in the 'correction' field should only contain letters`,
		},
	}

	for _, test := range testCases {
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()

			replacer := &misspell.Replacer{}

			err := appendExtraWords(replacer, test.extraWords)
			require.EqualError(t, err, test.expected)
		})
	}
}
