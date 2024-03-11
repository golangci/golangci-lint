package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConfiguration_Validate(t *testing.T) {
	testCases := []struct {
		desc string
		cfg  *Configuration
	}{
		{
			desc: "version",
			cfg: &Configuration{
				Version: "v1.57.0",
				Plugins: []*Plugin{
					{
						Module:  "example.org/foo/bar",
						Import:  "example.org/foo/bar/test",
						Version: "v1.2.3",
					},
				},
			},
		},
		{
			desc: "path",
			cfg: &Configuration{
				Version: "v1.57.0",
				Plugins: []*Plugin{
					{
						Module: "example.org/foo/bar",
						Import: "example.org/foo/bar/test",
						Path:   "/my/path",
					},
				},
			},
		},
	}

	for _, test := range testCases {
		test := test
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()

			err := test.cfg.Validate()
			require.NoError(t, err)
		})
	}
}

func TestConfiguration_Validate_error(t *testing.T) {
	testCases := []struct {
		desc     string
		cfg      *Configuration
		expected string
	}{
		{
			desc:     "missing version",
			cfg:      &Configuration{},
			expected: "root field 'version' is required",
		},
		{
			desc: "no plugins",
			cfg: &Configuration{
				Version: "v1.57.0",
			},
			expected: "no plugins defined",
		},
		{
			desc: "missing module",
			cfg: &Configuration{
				Version: "v1.57.0",
				Plugins: []*Plugin{
					{
						Module:  "",
						Import:  "example.org/foo/bar/test",
						Version: "v1.2.3",
						Path:    "/my/path",
					},
				},
			},
			expected: "field 'module' is required",
		},
		{
			desc: "module version and path",
			cfg: &Configuration{
				Version: "v1.57.0",
				Plugins: []*Plugin{
					{
						Module:  "example.org/foo/bar",
						Import:  "example.org/foo/bar/test",
						Version: "v1.2.3",
						Path:    "/my/path",
					},
				},
			},
			expected: "invalid configuration: 'version' and 'path' should not be provided at the same time",
		},
		{
			desc: "no module version and path",
			cfg: &Configuration{
				Version: "v1.57.0",
				Plugins: []*Plugin{
					{
						Module: "example.org/foo/bar",
						Import: "example.org/foo/bar/test",
					},
				},
			},
			expected: "missing information: 'version' or 'path' should be provided",
		},
	}

	for _, test := range testCases {
		test := test
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()

			err := test.cfg.Validate()

			assert.EqualError(t, err, test.expected)
		})
	}
}
