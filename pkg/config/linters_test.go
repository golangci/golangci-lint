package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLinters_validateDisabledAndEnabledAtOneMoment(t *testing.T) {
	testCases := []struct {
		desc string
		cfg  *Linters
	}{
		{
			desc: "2 different sets",
			cfg: &Linters{
				Enable:  []string{"dupl", "gofmt", "misspell"},
				Disable: []string{"goimports", "gosec", "nolintlint"},
			},
		},
		{
			desc: "only enable",
			cfg: &Linters{
				Enable:  []string{"goimports", "gosec", "nolintlint"},
				Disable: nil,
			},
		},
		{
			desc: "only disable",
			cfg: &Linters{
				Enable:  nil,
				Disable: []string{"dupl", "gofmt", "misspell"},
			},
		},
		{
			desc: "no sets",
			cfg: &Linters{
				Enable:  nil,
				Disable: nil,
			},
		},
	}

	for _, test := range testCases {
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()

			err := test.cfg.validateDisabledAndEnabledAtOneMoment()
			require.NoError(t, err)
		})
	}
}

func TestLinters_validateDisabledAndEnabledAtOneMoment_error(t *testing.T) {
	testCases := []struct {
		desc     string
		cfg      *Linters
		expected string
	}{
		{
			desc: "disable one linter of the enabled linters",
			cfg: &Linters{
				Enable:  []string{"dupl", "gofmt", "misspell"},
				Disable: []string{"dupl", "gosec", "nolintlint"},
			},
			expected: `linter "dupl" can't be disabled and enabled at one moment`,
		},
		{
			desc: "disable multiple enabled linters",
			cfg: &Linters{
				Enable:  []string{"dupl", "gofmt", "misspell"},
				Disable: []string{"dupl", "gofmt", "misspell"},
			},
			expected: `linter "dupl" can't be disabled and enabled at one moment`,
		},
	}

	for _, test := range testCases {
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()

			err := test.cfg.validateDisabledAndEnabledAtOneMoment()
			require.Error(t, err)

			require.EqualError(t, err, test.expected)
		})
	}
}

func TestLinters_validateAllDisableEnableOptions(t *testing.T) {
	testCases := []struct {
		desc string
		cfg  *Linters
	}{
		{
			desc: "nothing",
			cfg:  &Linters{},
		},
		{
			desc: "enable and disable",
			cfg: &Linters{
				Enable:     []string{"goimports", "gosec", "nolintlint"},
				EnableAll:  false,
				Disable:    []string{"dupl", "gofmt", "misspell"},
				DisableAll: false,
			},
		},
		{
			desc: "disable-all and enable",
			cfg: &Linters{
				Enable:     []string{"goimports", "gosec", "nolintlint"},
				EnableAll:  false,
				Disable:    nil,
				DisableAll: true,
			},
		},
		{
			desc: "enable-all and disable",
			cfg: &Linters{
				Enable:     nil,
				EnableAll:  true,
				Disable:    []string{"goimports", "gosec", "nolintlint"},
				DisableAll: false,
			},
		},
		{
			desc: "enable-all and enable and fast",
			cfg: &Linters{
				Enable:     []string{"dupl", "gofmt", "misspell"},
				EnableAll:  true,
				Disable:    nil,
				DisableAll: false,
				Fast:       true,
			},
		},
	}

	for _, test := range testCases {
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()

			err := test.cfg.validateAllDisableEnableOptions()
			require.NoError(t, err)
		})
	}
}

func TestLinters_validateAllDisableEnableOptions_error(t *testing.T) {
	testCases := []struct {
		desc     string
		cfg      *Linters
		expected string
	}{
		{
			desc: "enable-all and disable-all",
			cfg: &Linters{
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
			cfg: &Linters{
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
			cfg: &Linters{
				Enable:     []string{"nolintlint"},
				EnableAll:  false,
				Disable:    []string{"dupl", "gofmt", "misspell"},
				DisableAll: true,
				Fast:       false,
			},
			expected: "can't combine options --disable-all and --disable",
		},
		{
			desc: "enable-all and enable",
			cfg: &Linters{
				Enable:     []string{"dupl", "gofmt", "misspell"},
				EnableAll:  true,
				Disable:    nil,
				DisableAll: false,
				Fast:       false,
			},
			expected: "can't combine options --enable-all and --enable",
		},
	}

	for _, test := range testCases {
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()

			err := test.cfg.validateAllDisableEnableOptions()
			require.Error(t, err)

			require.EqualError(t, err, test.expected)
		})
	}
}
