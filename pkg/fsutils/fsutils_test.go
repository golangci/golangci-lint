package fsutils

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestShortestRelPath(t *testing.T) {
	testCases := []struct {
		desc     string
		path     string
		wd       string
		expected string
	}{
		{
			desc:     "based on parent path",
			path:     "fsutils_test.go",
			wd:       filepath.Join("..", "fsutils"),
			expected: "fsutils_test.go",
		},
		{
			desc:     "based on current working directory",
			path:     "fsutils_test.go",
			wd:       "",
			expected: "fsutils_test.go",
		},
	}

	for _, test := range testCases {
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()

			rel, err := ShortestRelPath("fsutils_test.go", filepath.Join("..", "fsutils"))
			require.NoError(t, err)

			assert.Equal(t, test.expected, rel)
		})
	}
}
