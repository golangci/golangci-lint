package goutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCheckGoVersion(t *testing.T) {
	testCases := []struct {
		desc    string
		version string
		require require.ErrorAssertionFunc
	}{
		{
			desc:    "version greater than runtime version (patch)",
			version: "1.30.1",
			require: require.Error,
		},
		{
			desc:    "version greater than runtime version (family)",
			version: "1.30",
			require: require.Error,
		},
		{
			desc:    "version greater than runtime version (RC)",
			version: "1.30.0-rc1",
			require: require.Error,
		},
		{
			desc: "version equals to runtime version",
			version: func() string {
				rv, _ := CleanRuntimeVersion()
				return rv
			}(),
			require: require.NoError,
		},
		{
			desc:    "version lower than runtime version (patch)",
			version: "1.19.1",
			require: require.NoError,
		},
		{
			desc:    "version lower than runtime version (family)",
			version: "1.19",
			require: require.NoError,
		},
		{
			desc:    "version lower than runtime version (RC)",
			version: "1.19.0-rc1",
			require: require.NoError,
		},
	}

	for _, test := range testCases {
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()

			err := CheckGoVersion(test.version)
			test.require(t, err)
		})
	}
}

func TestTrimGoVersion(t *testing.T) {
	testCases := []struct {
		desc     string
		version  string
		expected string
	}{
		{
			desc:     "patched version",
			version:  "1.22.0",
			expected: "1.22",
		},
		{
			desc:     "minor version",
			version:  "1.22",
			expected: "1.22",
		},
		{
			desc:     "RC version",
			version:  "1.22rc1",
			expected: "1.22",
		},
		{
			desc:     "alpha version",
			version:  "1.22alpha1",
			expected: "1.22",
		},
		{
			desc:     "beta version",
			version:  "1.22beta1",
			expected: "1.22",
		},
		{
			desc:     "semver RC version",
			version:  "1.22.0-rc1",
			expected: "1.22",
		},
	}

	for _, test := range testCases {
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()

			version := TrimGoVersion(test.version)
			assert.Equal(t, test.expected, version)
		})
	}
}

func Test_cleanRuntimeVersion(t *testing.T) {
	testCases := []struct {
		desc     string
		version  string
		expected string
	}{
		{
			desc:     "go version",
			version:  "go1.22.0",
			expected: "go1.22.0",
		},
		{
			desc:     "language version",
			version:  "go1.22",
			expected: "go1.22",
		},
		{
			desc:     "GOEXPERIMENT",
			version:  "go1.23.0 X:boringcrypto",
			expected: "go1.23.0",
		},
		{
			desc:     "devel",
			version:  "devel go1.24-e705a2d Wed Aug 7 01:16:42 2024 +0000 linux/amd64",
			expected: "go1.24-e705a2d",
		},
	}

	for _, test := range testCases {
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()

			v, err := cleanRuntimeVersion(test.version)
			require.NoError(t, err)

			assert.Equal(t, test.expected, v)
		})
	}
}
