package lintersdb

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_trimGoVersion(t *testing.T) {
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
		test := test
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()

			version := trimGoVersion(test.version)
			assert.Equal(t, test.expected, version)
		})
	}
}
