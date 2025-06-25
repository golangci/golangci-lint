package internal

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

func Test_createVersion(t *testing.T) {
	parse, err := time.Parse(time.RFC3339Nano, "2010-02-04T21:00:57.123456789+08:00")
	require.NoError(t, err)

	version := createVersion("1.2.3", parse)

	assert.Equal(t, "1.2.3-custom-gcl-1265288457123456789", version)
}
