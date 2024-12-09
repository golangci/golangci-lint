package config

import (
	"os/user"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_expandHomeDir(t *testing.T) {
	u, err := user.Current()
	require.NoError(t, err)

	testCases := []struct {
		path        string
		expected    string
		expectedErr bool
	}{
		{path: "", expected: ""},
		{path: "~", expected: u.HomeDir},
		{path: "/foo", expected: "/foo"},
		{path: "\\foo", expected: "\\foo"},
		{path: "C:\foo", expected: "C:\foo"},
		{path: "~/foo/bar", expected: filepath.Join(u.HomeDir, "foo", "bar")},
		{path: "~foo/foo", expectedErr: true},
	}

	for _, tc := range testCases {
		actual, err := expandHomeDir(tc.path)

		if tc.expectedErr {
			require.Error(t, err)
			return
		}
		require.NoError(t, err)
		assert.Equal(t, tc.expected, actual)
	}
}
