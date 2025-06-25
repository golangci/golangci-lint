package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_sanitizeVersion(t *testing.T) {
	testCases := []struct {
		desc     string
		input    string
		expected string
	}{
		{
			desc:     "ampersand",
			input:    " te&st",
			expected: "test",
		},
		{
			desc:     "pipe",
			input:    " te|st",
			expected: "test",
		},
		{
			desc:     "version",
			input:    "v1.2.3",
			expected: "v1.2.3",
		},
		{
			desc:     "branch",
			input:    "feat/test",
			expected: "feat/test",
		},
		{
			desc:     "branch",
			input:    "value --key",
			expected: "valuekey",
		},
		{
			desc:     "hash",
			input:    "cd8b1177",
			expected: "cd8b1177",
		},
	}

	for _, test := range testCases {
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()

			v := sanitizeVersion(test.input)

			assert.Equal(t, test.expected, v)
		})
	}
}
