package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_buildConfig(t *testing.T) {
	testCases := []struct {
		desc       string
		inputPath  string
		minVersion version
		expected   string
	}{
		{
			desc:       "v1",
			inputPath:  "all-releases.json",
			minVersion: version{major: 1, minor: 28, patch: 3},
			expected:   "github-action-config.json",
		},
		{
			desc:       "v1 only",
			inputPath:  "all-releases-v2.json",
			minVersion: version{major: 1, minor: 28, patch: 3},
			expected:   "github-action-config-v1.json",
		},
	}

	for _, test := range testCases {
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()

			allReleases := unmarshalRelease(t, test.inputPath)

			config, err := buildConfig(allReleases, test.minVersion)
			require.NoError(t, err)

			data, err := json.MarshalIndent(config, "", "  ")
			require.NoError(t, err)

			expected, err := os.ReadFile(filepath.Join("testdata", test.expected))
			require.NoError(t, err)

			assert.JSONEq(t, string(expected), string(data))
		})
	}
}

func unmarshalRelease(t *testing.T, filename string) []release {
	file, err := os.Open(filepath.Join("testdata", filename))
	require.NoError(t, err)

	t.Cleanup(func() {
		_ = file.Close()
	})

	var data []release
	err = json.NewDecoder(file).Decode(&data)
	require.NoError(t, err)

	return data
}
