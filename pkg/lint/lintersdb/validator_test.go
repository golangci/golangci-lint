package lintersdb

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/logutils"
)

type validateErrorTestCase struct {
	desc     string
	cfg      *config.Linters
	expected string
}

var validateLintersNamesErrorTestCases = []validateErrorTestCase{
	{
		desc: "unknown enabled linter",
		cfg: &config.Linters{
			Enable:  []string{"golangci"},
			Disable: nil,
		},
		expected: `unknown linters: 'golangci', run 'golangci-lint help linters' to see the list of supported linters`,
	},
	{
		desc: "unknown disabled linter",
		cfg: &config.Linters{
			Enable:  nil,
			Disable: []string{"golangci"},
		},
		expected: `unknown linters: 'golangci', run 'golangci-lint help linters' to see the list of supported linters`,
	},
}

var validatePresetsErrorTestCases = []validateErrorTestCase{
	{
		desc: "unknown preset",
		cfg: &config.Linters{
			EnableAll: false,
			Presets:   []string{"golangci"},
		},
		expected: "no such preset \"golangci\": only next presets exist: " +
			"(bugs|comment|complexity|error|format|import|metalinter|module|performance|sql|style|test|unused)",
	},
	{
		desc: "presets and enable-all",
		cfg: &config.Linters{
			EnableAll: true,
			Presets:   []string{"bugs"},
		},
		expected: `--presets is incompatible with --enable-all`,
	},
}

type validatorTestCase struct {
	desc string
	cfg  *config.Linters
}

var validateLintersNamesTestCases = []validatorTestCase{
	{
		desc: "no enable no disable",
		cfg: &config.Linters{
			Enable:  nil,
			Disable: nil,
		},
	},
	{
		desc: "existing enabled linter",
		cfg: &config.Linters{
			Enable:  []string{"gofmt"},
			Disable: nil,
		},
	},
	{
		desc: "existing disabled linter",
		cfg: &config.Linters{
			Enable:  nil,
			Disable: []string{"gofmt"},
		},
	},
}

var validatePresetsTestCases = []validatorTestCase{
	{
		desc: "known preset",
		cfg: &config.Linters{
			EnableAll: false,
			Presets:   []string{"bugs"},
		},
	},
	{
		desc: "enable-all and no presets",
		cfg: &config.Linters{
			EnableAll: true,
			Presets:   nil,
		},
	},
	{
		desc: "no presets",
		cfg: &config.Linters{
			EnableAll: false,
			Presets:   nil,
		},
	},
}

func TestValidator_Validate(t *testing.T) {
	m, err := NewManager(nil, nil, NewLinterBuilder())
	require.NoError(t, err)

	v := NewValidator(m)

	var testCases []validatorTestCase
	testCases = append(testCases, validateLintersNamesTestCases...)
	testCases = append(testCases, validatePresetsTestCases...)

	for _, test := range testCases {
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()

			err := v.Validate(&config.Config{Linters: *test.cfg})
			require.NoError(t, err)
		})
	}
}

func TestValidator_Validate_error(t *testing.T) {
	m, err := NewManager(nil, nil, NewLinterBuilder())
	require.NoError(t, err)

	v := NewValidator(m)

	var testCases []validateErrorTestCase
	testCases = append(testCases, validateLintersNamesErrorTestCases...)
	testCases = append(testCases, validatePresetsErrorTestCases...)

	for _, test := range testCases {
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()

			err := v.Validate(&config.Config{Linters: *test.cfg})
			require.Error(t, err)

			require.EqualError(t, err, test.expected)
		})
	}
}

func TestValidator_validateLintersNames(t *testing.T) {
	m, err := NewManager(nil, nil, NewLinterBuilder())
	require.NoError(t, err)

	v := NewValidator(m)

	for _, test := range validateLintersNamesTestCases {
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()

			err := v.validateLintersNames(test.cfg)
			require.NoError(t, err)
		})
	}
}

func TestValidator_validateLintersNames_error(t *testing.T) {
	m, err := NewManager(nil, nil, NewLinterBuilder())
	require.NoError(t, err)

	v := NewValidator(m)

	for _, test := range validateLintersNamesErrorTestCases {
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()

			err := v.validateLintersNames(test.cfg)
			require.Error(t, err)

			require.EqualError(t, err, test.expected)
		})
	}
}

func TestValidator_validatePresets(t *testing.T) {
	v := NewValidator(nil)

	for _, test := range validatePresetsTestCases {
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()

			err := v.validatePresets(test.cfg)
			require.NoError(t, err)
		})
	}
}

func TestValidator_validatePresets_error(t *testing.T) {
	v := NewValidator(nil)

	for _, test := range validatePresetsErrorTestCases {
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()

			err := v.validatePresets(test.cfg)
			require.Error(t, err)

			require.EqualError(t, err, test.expected)
		})
	}
}

func TestValidator_alternativeNamesDeprecation(t *testing.T) {
	t.Setenv(logutils.EnvTestRun, "0")

	log := logutils.NewMockLog().
		OnWarnf("The name %q is deprecated. The linter has been renamed to: %s.", "vet", "govet").
		OnWarnf("The name %q is deprecated. The linter has been renamed to: %s.", "vetshadow", "govet").
		OnWarnf("The name %q is deprecated. The linter has been renamed to: %s.", "logrlint", "loggercheck").
		OnWarnf("The linter named %q is deprecated. It has been split into: %s.", "megacheck", "gosimple, staticcheck, unused").
		OnWarnf("The name %q is deprecated. The linter has been renamed to: %s.", "gas", "gosec")

	m, err := NewManager(log, nil, NewLinterBuilder())
	require.NoError(t, err)

	v := NewValidator(m)

	cfg := &config.Linters{
		Enable:  []string{"vet", "vetshadow", "logrlint"},
		Disable: []string{"megacheck", "gas"},
	}

	err = v.alternativeNamesDeprecation(cfg)
	require.NoError(t, err)
}
