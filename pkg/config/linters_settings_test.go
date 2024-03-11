package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLintersSettings_Validate(t *testing.T) {
	testCases := []struct {
		desc     string
		settings *LintersSettings
	}{
		{
			desc: "custom linter",
			settings: &LintersSettings{
				Custom: map[string]CustomLinterSettings{
					"example": {
						Type: "module",
					},
				},
			},
		},
		{
			desc: "govet",
			settings: &LintersSettings{
				Govet: GovetSettings{
					Enable:     []string{"a"},
					DisableAll: true,
				},
			},
		},
	}

	for _, test := range testCases {
		test := test
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()

			err := test.settings.Validate()
			assert.NoError(t, err)
		})
	}
}

func TestLintersSettings_Validate_error(t *testing.T) {
	testCases := []struct {
		desc     string
		settings *LintersSettings
		expected string
	}{
		{
			desc: "custom linter error",
			settings: &LintersSettings{
				Custom: map[string]CustomLinterSettings{
					"example": {
						Type: "module",
						Path: "example",
					},
				},
			},
			expected: `custom linter "example": path not supported with module type`,
		},
		{
			desc: "govet error",
			settings: &LintersSettings{
				Govet: GovetSettings{
					EnableAll:  true,
					DisableAll: true,
				},
			},
			expected: "govet: enable-all and disable-all can't be combined",
		},
	}

	for _, test := range testCases {
		test := test
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()

			err := test.settings.Validate()

			assert.EqualError(t, err, test.expected)
		})
	}
}

func TestCustomLinterSettings_Validate(t *testing.T) {
	testCases := []struct {
		desc     string
		settings *CustomLinterSettings
	}{
		{
			desc: "only path",
			settings: &CustomLinterSettings{
				Path: "example",
			},
		},
		{
			desc: "path and type goplugin",
			settings: &CustomLinterSettings{
				Type: "goplugin",
				Path: "example",
			},
		},
		{
			desc: "type module",
			settings: &CustomLinterSettings{
				Type: "module",
			},
		},
	}

	for _, test := range testCases {
		test := test
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()

			err := test.settings.Validate()
			assert.NoError(t, err)
		})
	}
}

func TestCustomLinterSettings_Validate_error(t *testing.T) {
	testCases := []struct {
		desc     string
		settings *CustomLinterSettings
		expected string
	}{
		{
			desc:     "missing path",
			settings: &CustomLinterSettings{},
			expected: "path is required",
		},
		{
			desc: "module and path",
			settings: &CustomLinterSettings{
				Type: "module",
				Path: "example",
			},
			expected: "path not supported with module type",
		},
	}

	for _, test := range testCases {
		test := test
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()

			err := test.settings.Validate()

			assert.EqualError(t, err, test.expected)
		})
	}
}

func TestGovetSettings_Validate(t *testing.T) {
	testCases := []struct {
		desc     string
		settings *GovetSettings
	}{
		{
			desc:     "empty",
			settings: &GovetSettings{},
		},
		{
			desc: "disable-all and enable",
			settings: &GovetSettings{
				Enable:     []string{"a"},
				DisableAll: true,
			},
		},
		{
			desc: "enable-all and disable",
			settings: &GovetSettings{
				Disable:   []string{"a"},
				EnableAll: true,
			},
		},
	}

	for _, test := range testCases {
		test := test
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()

			err := test.settings.Validate()
			assert.NoError(t, err)
		})
	}
}

func TestGovetSettings_Validate_error(t *testing.T) {
	testCases := []struct {
		desc     string
		settings *GovetSettings
		expected string
	}{
		{
			desc: "enable-all and disable-all",
			settings: &GovetSettings{
				EnableAll:  true,
				DisableAll: true,
			},
			expected: "govet: enable-all and disable-all can't be combined",
		},
		{
			desc: "enable-all and enable",
			settings: &GovetSettings{
				EnableAll: true,
				Enable:    []string{"a"},
			},
			expected: "govet: enable-all and enable can't be combined",
		},
		{
			desc: "disable-all and disable",
			settings: &GovetSettings{
				DisableAll: true,
				Disable:    []string{"a"},
			},
			expected: "govet: disable-all and disable can't be combined",
		},
	}

	for _, test := range testCases {
		test := test
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()

			err := test.settings.Validate()

			assert.EqualError(t, err, test.expected)
		})
	}
}
