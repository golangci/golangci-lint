package lintersdb

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/golangci/golangci-lint/pkg/config"
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

var validateDisabledAndEnabledAtOneMomentErrorTestCases = []validateErrorTestCase{
	{
		desc: "disable one linter of the enabled linters",
		cfg: &config.Linters{
			Enable:  []string{"dupl", "gofmt", "misspell"},
			Disable: []string{"dupl", "gosec", "nolintlint"},
		},
		expected: `linter "dupl" can't be disabled and enabled at one moment`,
	},
	{
		desc: "disable multiple enabled linters",
		cfg: &config.Linters{
			Enable:  []string{"dupl", "gofmt", "misspell"},
			Disable: []string{"dupl", "gofmt", "misspell"},
		},
		expected: `linter "dupl" can't be disabled and enabled at one moment`,
	},
}

var validateAllDisableEnableOptionsErrorTestCases = []validateErrorTestCase{
	{
		desc: "enable-all and disable-all",
		cfg: &config.Linters{
			Enable:     nil,
			EnableAll:  true,
			Disable:    nil,
			DisableAll: true,
			Fast:       false,
		},
		expected: "--enable-all and --disable-all options must not be combined",
	},
	{
		desc: "disable-all and disable no enable no preset",
		cfg: &config.Linters{
			Enable:     nil,
			EnableAll:  false,
			Disable:    []string{"dupl", "gofmt", "misspell"},
			DisableAll: true,
			Fast:       false,
		},
		expected: "all linters were disabled, but no one linter was enabled: must enable at least one",
	},
	{
		desc: "disable-all and disable with enable",
		cfg: &config.Linters{
			Enable:     []string{"nolintlint"},
			EnableAll:  false,
			Disable:    []string{"dupl", "gofmt", "misspell"},
			DisableAll: true,
			Fast:       false,
		},
		expected: "can't combine options --disable-all and --disable dupl",
	},
	{
		desc: "enable-all and enable",
		cfg: &config.Linters{
			Enable:     []string{"dupl", "gofmt", "misspell"},
			EnableAll:  true,
			Disable:    nil,
			DisableAll: false,
			Fast:       false,
		},
		expected: "can't combine options --enable-all and --enable dupl",
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

var validateDisabledAndEnabledAtOneMomentTestCases = []validatorTestCase{
	{
		desc: "2 different sets",
		cfg: &config.Linters{
			Enable:  []string{"dupl", "gofmt", "misspell"},
			Disable: []string{"goimports", "gosec", "nolintlint"},
		},
	},
	{
		desc: "only enable",
		cfg: &config.Linters{
			Enable:  []string{"goimports", "gosec", "nolintlint"},
			Disable: nil,
		},
	},
	{
		desc: "only disable",
		cfg: &config.Linters{
			Enable:  nil,
			Disable: []string{"dupl", "gofmt", "misspell"},
		},
	},
	{
		desc: "no sets",
		cfg: &config.Linters{
			Enable:  nil,
			Disable: nil,
		},
	},
}

var validateAllDisableEnableOptionsTestCases = []validatorTestCase{
	{
		desc: "nothing",
		cfg:  &config.Linters{},
	},
	{
		desc: "enable and disable",
		cfg: &config.Linters{
			Enable:     []string{"goimports", "gosec", "nolintlint"},
			EnableAll:  false,
			Disable:    []string{"dupl", "gofmt", "misspell"},
			DisableAll: false,
		},
	},
	{
		desc: "disable-all and enable",
		cfg: &config.Linters{
			Enable:     []string{"goimports", "gosec", "nolintlint"},
			EnableAll:  false,
			Disable:    nil,
			DisableAll: true,
		},
	},
	{
		desc: "enable-all and disable",
		cfg: &config.Linters{
			Enable:     nil,
			EnableAll:  true,
			Disable:    []string{"goimports", "gosec", "nolintlint"},
			DisableAll: false,
		},
	},
	{
		desc: "enable-all and enable and fast",
		cfg: &config.Linters{
			Enable:     []string{"dupl", "gofmt", "misspell"},
			EnableAll:  true,
			Disable:    nil,
			DisableAll: false,
			Fast:       true,
		},
	},
}

func TestValidator_validateEnabledDisabledLintersConfig(t *testing.T) {
	v := NewValidator(NewManager(nil, nil))

	var testCases []validatorTestCase
	testCases = append(testCases, validateLintersNamesTestCases...)
	testCases = append(testCases, validatePresetsTestCases...)
	testCases = append(testCases, validateDisabledAndEnabledAtOneMomentTestCases...)
	testCases = append(testCases, validateAllDisableEnableOptionsTestCases...)

	for _, test := range testCases {
		test := test
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()

			err := v.validateEnabledDisabledLintersConfig(test.cfg)
			require.NoError(t, err)
		})
	}
}

func TestValidator_validateEnabledDisabledLintersConfig_error(t *testing.T) {
	v := NewValidator(NewManager(nil, nil))

	var testCases []validateErrorTestCase
	testCases = append(testCases, validateLintersNamesErrorTestCases...)
	testCases = append(testCases, validatePresetsErrorTestCases...)
	testCases = append(testCases, validateDisabledAndEnabledAtOneMomentErrorTestCases...)
	testCases = append(testCases, validateAllDisableEnableOptionsErrorTestCases...)

	for _, test := range testCases {
		test := test
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()

			err := v.validateEnabledDisabledLintersConfig(test.cfg)
			require.Error(t, err)

			require.EqualError(t, err, test.expected)
		})
	}
}

func TestValidator_validateLintersNames(t *testing.T) {
	v := NewValidator(NewManager(nil, nil))

	for _, test := range validateLintersNamesTestCases {
		test := test
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()

			err := v.validateLintersNames(test.cfg)
			require.NoError(t, err)
		})
	}
}

func TestValidator_validateLintersNames_error(t *testing.T) {
	v := NewValidator(NewManager(nil, nil))

	for _, test := range validateLintersNamesErrorTestCases {
		test := test
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
		test := test
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
		test := test
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()

			err := v.validatePresets(test.cfg)
			require.Error(t, err)

			require.EqualError(t, err, test.expected)
		})
	}
}

func TestValidator_validateDisabledAndEnabledAtOneMoment(t *testing.T) {
	v := NewValidator(nil)

	for _, test := range validateDisabledAndEnabledAtOneMomentTestCases {
		test := test
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()

			err := v.validateDisabledAndEnabledAtOneMoment(test.cfg)
			require.NoError(t, err)
		})
	}
}

func TestValidator_validateDisabledAndEnabledAtOneMoment_error(t *testing.T) {
	v := NewValidator(nil)

	for _, test := range validateDisabledAndEnabledAtOneMomentErrorTestCases {
		test := test
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()

			err := v.validateDisabledAndEnabledAtOneMoment(test.cfg)
			require.Error(t, err)

			require.EqualError(t, err, test.expected)
		})
	}
}

func TestValidator_validateAllDisableEnableOptions(t *testing.T) {
	v := NewValidator(nil)

	for _, test := range validateAllDisableEnableOptionsTestCases {
		test := test
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()

			err := v.validateAllDisableEnableOptions(test.cfg)
			require.NoError(t, err)
		})
	}
}

func TestValidator_validateAllDisableEnableOptions_error(t *testing.T) {
	v := NewValidator(nil)

	for _, test := range validateAllDisableEnableOptionsErrorTestCases {
		test := test
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()

			err := v.validateAllDisableEnableOptions(test.cfg)
			require.Error(t, err)

			require.EqualError(t, err, test.expected)
		})
	}
}
