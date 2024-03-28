package processors

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_absDirs(t *testing.T) {
	testCases := []struct {
		desc     string
		args     []string
		expected []string
	}{
		{
			desc:     "empty",
			expected: []string{mustAbs(t, ".")},
		},
		{
			desc:     "wildcard",
			args:     []string{"./..."},
			expected: []string{mustAbs(t, ".")},
		},
		{
			desc:     "wildcard directory",
			args:     []string{"foo/..."},
			expected: []string{mustAbs(t, "foo")},
		},
		{
			desc:     "Go file",
			args:     []string{"./foo/bar.go"},
			expected: []string{mustAbs(t, "foo")},
		},
		{
			desc:     "relative directory",
			args:     []string{filepath.FromSlash("./foo")},
			expected: []string{mustAbs(t, "foo")},
		},
		{
			desc:     "absolute directory",
			args:     []string{mustAbs(t, "foo")},
			expected: []string{mustAbs(t, "foo")},
		},
	}

	for _, test := range testCases {
		test := test
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()

			results, err := absDirs(test.args)
			require.NoError(t, err)

			assert.Equal(t, test.expected, results)
		})
	}
}

func mustAbs(t *testing.T, p string) string {
	t.Helper()

	abs, err := filepath.Abs(p)
	require.NoError(t, err)

	return abs
}
