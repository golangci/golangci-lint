package lintersdb

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/golangci/golangci-lint/v2/pkg/config"
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

func TestValidator_Validate(t *testing.T) {
	m, err := NewManager(nil, nil, NewLinterBuilder())
	require.NoError(t, err)

	v := NewValidator(m)

	var testCases []validatorTestCase
	testCases = append(testCases, validateLintersNamesTestCases...)

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
