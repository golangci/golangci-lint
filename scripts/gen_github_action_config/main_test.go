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
	file, err := os.Open(filepath.Join("testdata", "all-releases.json"))
	require.NoError(t, err)

	defer file.Close()

	var allReleases []release
	err = json.NewDecoder(file).Decode(&allReleases)
	require.NoError(t, err)

	minAllowedVersion := version{major: 1, minor: 51, patch: 0}

	config, err := buildConfig(allReleases, minAllowedVersion)
	require.NoError(t, err)

	expected := &actionConfig{MinorVersionToConfig: map[string]versionConfig{
		"latest": {Error: "", TargetVersion: "v1.63.4", AssetURL: "https://github.com/golangci/golangci-lint/releases/download/v1.63.4/golangci-lint-1.63.4-linux-amd64.tar.gz"},
		"v1.3":   {Error: "golangci-lint version 'v1.3' isn't supported: we support only v1.51.0 and later versions"},
		"v1.4":   {Error: "golangci-lint version 'v1.4' isn't supported: we support only v1.51.0 and later versions"},
		"v1.5":   {Error: "golangci-lint version 'v1.5' isn't supported: we support only v1.51.0 and later versions"},
		"v1.6":   {Error: "golangci-lint version 'v1.6' isn't supported: we support only v1.51.0 and later versions"},
		"v1.7":   {Error: "golangci-lint version 'v1.7' isn't supported: we support only v1.51.0 and later versions"},
		"v1.8":   {Error: "golangci-lint version 'v1.8' isn't supported: we support only v1.51.0 and later versions"},
		"v1.9":   {Error: "golangci-lint version 'v1.9' isn't supported: we support only v1.51.0 and later versions"},
		"v1.10":  {Error: "golangci-lint version 'v1.10' isn't supported: we support only v1.51.0 and later versions"},
		"v1.11":  {Error: "golangci-lint version 'v1.11' isn't supported: we support only v1.51.0 and later versions"},
		"v1.12":  {Error: "golangci-lint version 'v1.12' isn't supported: we support only v1.51.0 and later versions"},
		"v1.13":  {Error: "golangci-lint version 'v1.13' isn't supported: we support only v1.51.0 and later versions"},
		"v1.14":  {Error: "golangci-lint version 'v1.14' isn't supported: we support only v1.51.0 and later versions"},
		"v1.15":  {Error: "golangci-lint version 'v1.15' isn't supported: we support only v1.51.0 and later versions"},
		"v1.16":  {Error: "golangci-lint version 'v1.16' isn't supported: we support only v1.51.0 and later versions"},
		"v1.17":  {Error: "golangci-lint version 'v1.17' isn't supported: we support only v1.51.0 and later versions"},
		"v1.18":  {Error: "golangci-lint version 'v1.18' isn't supported: we support only v1.51.0 and later versions"},
		"v1.19":  {Error: "golangci-lint version 'v1.19' isn't supported: we support only v1.51.0 and later versions"},
		"v1.20":  {Error: "golangci-lint version 'v1.20' isn't supported: we support only v1.51.0 and later versions"},
		"v1.21":  {Error: "golangci-lint version 'v1.21' isn't supported: we support only v1.51.0 and later versions"},
		"v1.22":  {Error: "golangci-lint version 'v1.22' isn't supported: we support only v1.51.0 and later versions"},
		"v1.23":  {Error: "golangci-lint version 'v1.23' isn't supported: we support only v1.51.0 and later versions"},
		"v1.24":  {Error: "golangci-lint version 'v1.24' isn't supported: we support only v1.51.0 and later versions"},
		"v1.25":  {Error: "golangci-lint version 'v1.25' isn't supported: we support only v1.51.0 and later versions"},
		"v1.26":  {Error: "golangci-lint version 'v1.26' isn't supported: we support only v1.51.0 and later versions"},
		"v1.27":  {Error: "golangci-lint version 'v1.27' isn't supported: we support only v1.51.0 and later versions"},
		"v1.28":  {Error: "golangci-lint version 'v1.28' isn't supported: we support only v1.51.0 and later versions"},
		"v1.29":  {Error: "golangci-lint version 'v1.29' isn't supported: we support only v1.51.0 and later versions"},
		"v1.30":  {Error: "golangci-lint version 'v1.30' isn't supported: we support only v1.51.0 and later versions"},
		"v1.31":  {Error: "golangci-lint version 'v1.31' isn't supported: we support only v1.51.0 and later versions"},
		"v1.32":  {Error: "golangci-lint version 'v1.32' isn't supported: we support only v1.51.0 and later versions"},
		"v1.33":  {Error: "golangci-lint version 'v1.33' isn't supported: we support only v1.51.0 and later versions"},
		"v1.34":  {Error: "golangci-lint version 'v1.34' isn't supported: we support only v1.51.0 and later versions"},
		"v1.35":  {Error: "golangci-lint version 'v1.35' isn't supported: we support only v1.51.0 and later versions"},
		"v1.36":  {Error: "golangci-lint version 'v1.36' isn't supported: we support only v1.51.0 and later versions"},
		"v1.37":  {Error: "golangci-lint version 'v1.37' isn't supported: we support only v1.51.0 and later versions"},
		"v1.38":  {Error: "golangci-lint version 'v1.38' isn't supported: we support only v1.51.0 and later versions"},
		"v1.39":  {Error: "golangci-lint version 'v1.39' isn't supported: we support only v1.51.0 and later versions"},
		"v1.40":  {Error: "golangci-lint version 'v1.40' isn't supported: we support only v1.51.0 and later versions"},
		"v1.41":  {Error: "golangci-lint version 'v1.41' isn't supported: we support only v1.51.0 and later versions"},
		"v1.42":  {Error: "golangci-lint version 'v1.42' isn't supported: we support only v1.51.0 and later versions"},
		"v1.43":  {Error: "golangci-lint version 'v1.43' isn't supported: we support only v1.51.0 and later versions"},
		"v1.44":  {Error: "golangci-lint version 'v1.44' isn't supported: we support only v1.51.0 and later versions"},
		"v1.45":  {Error: "golangci-lint version 'v1.45' isn't supported: we support only v1.51.0 and later versions"},
		"v1.46":  {Error: "golangci-lint version 'v1.46' isn't supported: we support only v1.51.0 and later versions"},
		"v1.47":  {Error: "golangci-lint version 'v1.47' isn't supported: we support only v1.51.0 and later versions"},
		"v1.48":  {Error: "golangci-lint version 'v1.48' isn't supported: we support only v1.51.0 and later versions"},
		"v1.49":  {Error: "golangci-lint version 'v1.49' isn't supported: we support only v1.51.0 and later versions"},
		"v1.50":  {Error: "golangci-lint version 'v1.50' isn't supported: we support only v1.51.0 and later versions"},
		"v1.51":  {Error: "", TargetVersion: "v1.51.2", AssetURL: "https://github.com/golangci/golangci-lint/releases/download/v1.51.2/golangci-lint-1.51.2-linux-amd64.tar.gz"},
		"v1.52":  {Error: "", TargetVersion: "v1.52.2", AssetURL: "https://github.com/golangci/golangci-lint/releases/download/v1.52.2/golangci-lint-1.52.2-linux-amd64.tar.gz"},
		"v1.53":  {Error: "", TargetVersion: "v1.53.3", AssetURL: "https://github.com/golangci/golangci-lint/releases/download/v1.53.3/golangci-lint-1.53.3-linux-amd64.tar.gz"},
		"v1.54":  {Error: "", TargetVersion: "v1.54.2", AssetURL: "https://github.com/golangci/golangci-lint/releases/download/v1.54.2/golangci-lint-1.54.2-linux-amd64.tar.gz"},
		"v1.55":  {Error: "", TargetVersion: "v1.55.2", AssetURL: "https://github.com/golangci/golangci-lint/releases/download/v1.55.2/golangci-lint-1.55.2-linux-amd64.tar.gz"},
		"v1.56":  {Error: "", TargetVersion: "v1.56.2", AssetURL: "https://github.com/golangci/golangci-lint/releases/download/v1.56.2/golangci-lint-1.56.2-linux-amd64.tar.gz"},
		"v1.57":  {Error: "", TargetVersion: "v1.57.2", AssetURL: "https://github.com/golangci/golangci-lint/releases/download/v1.57.2/golangci-lint-1.57.2-linux-amd64.tar.gz"},
		"v1.58":  {Error: "", TargetVersion: "v1.58.2", AssetURL: "https://github.com/golangci/golangci-lint/releases/download/v1.58.2/golangci-lint-1.58.2-linux-amd64.tar.gz"},
		"v1.59":  {Error: "", TargetVersion: "v1.59.1", AssetURL: "https://github.com/golangci/golangci-lint/releases/download/v1.59.1/golangci-lint-1.59.1-linux-amd64.tar.gz"},
		"v1.60":  {Error: "", TargetVersion: "v1.60.3", AssetURL: "https://github.com/golangci/golangci-lint/releases/download/v1.60.3/golangci-lint-1.60.3-linux-amd64.tar.gz"},
		"v1.61":  {Error: "", TargetVersion: "v1.61.0", AssetURL: "https://github.com/golangci/golangci-lint/releases/download/v1.61.0/golangci-lint-1.61.0-linux-amd64.tar.gz"},
		"v1.62":  {Error: "", TargetVersion: "v1.62.2", AssetURL: "https://github.com/golangci/golangci-lint/releases/download/v1.62.2/golangci-lint-1.62.2-linux-amd64.tar.gz"},
		"v1.63":  {Error: "", TargetVersion: "v1.63.4", AssetURL: "https://github.com/golangci/golangci-lint/releases/download/v1.63.4/golangci-lint-1.63.4-linux-amd64.tar.gz"},
	}}

	assert.Equal(t, expected, config)
}
