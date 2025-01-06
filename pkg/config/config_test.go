package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsGoGreaterThanOrEqual(t *testing.T) {
	testCases := []struct {
		desc    string
		current string
		limit   string
		assert  assert.BoolAssertionFunc
	}{
		{
			desc:    "current (with minor.major) lower than limit",
			current: "go1.21",
			limit:   "1.22",
			assert:  assert.False,
		},
		{
			desc:    "current (with 0 patch) lower than limit",
			current: "go1.21.0",
			limit:   "1.22",
			assert:  assert.False,
		},
		{
			desc:    "current (current with multiple patches) lower than limit",
			current: "go1.21.6",
			limit:   "1.22",
			assert:  assert.False,
		},
		{
			desc:    "current lower than limit (with minor.major)",
			current: "go1.22",
			limit:   "1.22",
			assert:  assert.True,
		},
		{
			desc:    "current lower than limit (with 0 patch)",
			current: "go1.22.0",
			limit:   "1.22",
			assert:  assert.True,
		},
		{
			desc:    "current lower than limit (current with multiple patches)",
			current: "go1.22.6",
			limit:   "1.22",
			assert:  assert.True,
		},
		{
			desc:    "current greater than limit",
			current: "go1.23.0",
			limit:   "1.22",
			assert:  assert.True,
		},
		{
			desc:    "current with no prefix",
			current: "1.22",
			limit:   "1.22",
			assert:  assert.True,
		},
		{
			desc:    "invalid current value",
			current: "go",
			limit:   "1.22",
			assert:  assert.False,
		},
		{
			desc:    "invalid limit value",
			current: "go1.22",
			limit:   "go",
			assert:  assert.False,
		},
	}

	for _, test := range testCases {
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()

			test.assert(t, IsGoGreaterThanOrEqual(test.current, test.limit))
		})
	}
}

func Test_parseGoVersion(t *testing.T) {
	testCases := []struct {
		desc     string
		version  string
		expected string
	}{
		{
			desc:     "empty version",
			version:  "",
			expected: "",
		},
		{
			desc:     "no prefixed version",
			version:  "1.23.0",
			expected: "1.23.0",
		},
		{
			desc:     "semver version",
			version:  "go1.23.0",
			expected: "1.23.0",
		},
		{
			desc:     "family version",
			version:  "go1.23",
			expected: "1.23",
		},
		{
			desc:     "prerelease version",
			version:  "go1.24rc1",
			expected: "1.24",
		},
	}

	for _, test := range testCases {
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()

			v := parseGoVersion(test.version)

			assert.Equal(t, test.expected, v)
		})
	}
}
