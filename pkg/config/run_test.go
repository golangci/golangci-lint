package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRun_Validate(t *testing.T) {
	testCases := []struct {
		desc     string
		settings *Run
	}{
		{
			desc: "modules-download-mode: mod",
			settings: &Run{
				ModulesDownloadMode: "mod",
			},
		},
		{
			desc: "modules-download-mode: readonly",
			settings: &Run{
				ModulesDownloadMode: "readonly",
			},
		},
		{
			desc: "modules-download-mode: vendor",
			settings: &Run{
				ModulesDownloadMode: "vendor",
			},
		},
		{
			desc: "modules-download-mode: empty",
			settings: &Run{
				ModulesDownloadMode: "",
			},
		},
	}

	for _, test := range testCases {
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()

			err := test.settings.Validate()
			require.NoError(t, err)
		})
	}
}

func TestRun_Validate_error(t *testing.T) {
	testCases := []struct {
		desc     string
		settings *Run
		expected string
	}{
		{
			desc: "modules-download-mode: invalid",
			settings: &Run{
				ModulesDownloadMode: "invalid",
			},
			expected: "invalid modules download path invalid, only (mod|readonly|vendor) allowed",
		},
	}

	for _, test := range testCases {
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()

			err := test.settings.Validate()
			require.EqualError(t, err, test.expected)
		})
	}
}
