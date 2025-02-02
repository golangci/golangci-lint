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
	allReleases := unmarshalRelease(t, "all-releases.json")

	minAllowedVersion := version{major: 1, minor: 28, patch: 3}

	config, err := buildConfig(allReleases, minAllowedVersion)
	require.NoError(t, err)

	data, err := json.MarshalIndent(config, "", "  ")
	require.NoError(t, err)

	expected, err := os.ReadFile(filepath.Join("testdata", "github-action-config.json"))
	require.NoError(t, err)

	assert.JSONEq(t, string(expected), string(data))
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
