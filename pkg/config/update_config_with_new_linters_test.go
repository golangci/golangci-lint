package config

import (
	"io"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/golangci/golangci-lint/pkg/logutils"
)

func TestUpdateConfigWithNewLinters(t *testing.T) {
	writeConfig := func(t *testing.T, cfg string) string {
		tmpFile, err := ioutil.TempFile("", "golangci-lint-test-config-*.yml")
		require.NoError(t, err)

		_, err = io.WriteString(tmpFile, cfg)
		require.NoError(t, err)

		err = tmpFile.Close()
		require.NoError(t, err)

		t.Cleanup(func() {
			os.Remove(tmpFile.Name())
		})
		return tmpFile.Name()
	}

	readConfig := func(t *testing.T, filePath string) Config {
		var cfg Config
		cmdLineCfg := Config{Run: Run{Config: filePath}}
		r := NewFileReader(&cfg, &cmdLineCfg, logutils.NewStderrLog("testing"))
		err := r.Read()
		require.NoError(t, err)
		return cfg
	}

	t.Run(`when the "linters" -> "enable" node exists, we add to it, matching the indent size`, func(t *testing.T) {
		cfgFilePath := writeConfig(t, `
linters:
   enable:
   - other-linter
`)
		err := UpdateConfigFileWithNewLinters(cfgFilePath, []string{"new-linter"})
		require.NoError(t, err)
		cfg := readConfig(t, cfgFilePath)
		require.Contains(t, cfg.Linters.Enable, "new-linter")

		data, err := ioutil.ReadFile(cfgFilePath)
		require.NoError(t, err)
		assert.Contains(t, string(data), "\n"+strings.Repeat(" ", 6)+"- new-linter",
			"indent size does not match")
	})

	t.Run(`when there is no "enable" node, we create one`, func(t *testing.T) {
		cfgFilePath := writeConfig(t, `
linters: {}
`)
		err := UpdateConfigFileWithNewLinters(cfgFilePath, []string{"new-linter"})
		require.NoError(t, err)
		cfg := readConfig(t, cfgFilePath)
		assert.Contains(t, cfg.Linters.Enable, "new-linter")
	})

	t.Run(`when the file is empty, we create values from scratch`, func(t *testing.T) {
		cfgFilePath := writeConfig(t, `
{}
`)
		err := UpdateConfigFileWithNewLinters(cfgFilePath, []string{"new-linter"})
		require.NoError(t, err)
		cfg := readConfig(t, cfgFilePath)
		assert.Contains(t, cfg.Linters.Enable, "new-linter")
	})

	t.Run(`when there is no "linters" node, we create one`, func(t *testing.T) {
		cfgFilePath := writeConfig(t, "")
		err := UpdateConfigFileWithNewLinters(cfgFilePath, []string{"new-linter"})
		require.NoError(t, err)
		cfg := readConfig(t, cfgFilePath)
		assert.Contains(t, cfg.Linters.Enable, "new-linter")
	})
}
