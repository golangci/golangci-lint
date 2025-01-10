package commands

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_formatDescription(t *testing.T) {
	testCases := []struct {
		desc     string
		doc      string
		expected string
	}{
		{
			desc:     "empty description",
			doc:      "",
			expected: "",
		},
		{
			desc:     "simple description",
			doc:      "this is a test",
			expected: "This is a test.",
		},
		{
			desc:     "formatted description",
			doc:      "This is a test.",
			expected: "This is a test.",
		},
		{
			desc:     "multiline description",
			doc:      "this is a test\nanother line\n",
			expected: "This is a test.",
		},
		{
			desc:     "leading newline",
			doc:      "\nThis is a test.",
			expected: "This is a test.",
		},
	}

	for _, test := range testCases {
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()

			v := formatDescription(test.doc)

			assert.Equal(t, test.expected, v)
		})
	}
}
