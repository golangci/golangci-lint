package lint

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_buildArgs(t *testing.T) {
	testCases := []struct {
		desc     string
		args     []string
		expected []string
	}{
		{
			desc:     "empty",
			args:     nil,
			expected: []string{"./..."},
		},
		{
			desc:     "start with a dot",
			args:     []string{filepath.FromSlash("./foo")},
			expected: []string{filepath.FromSlash("./foo")},
		},
		{
			desc:     "start without a dot",
			args:     []string{"foo"},
			expected: []string{filepath.FromSlash("./foo")},
		},
		{
			desc:     "absolute path",
			args:     []string{filepath.FromSlash("/tmp/foo")},
			expected: []string{filepath.FromSlash("/tmp/foo")},
		},
	}

	for _, test := range testCases {
		test := test
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()

			results := buildArgs(test.args)

			assert.Equal(t, test.expected, results)
		})
	}
}
