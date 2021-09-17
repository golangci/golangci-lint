package config

import (
	"path/filepath"
	"testing"

	"github.com/golangci/golangci-lint/pkg/logutils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestFileReader_Read(t *testing.T) {
	cases := map[string]Config{
		"preset_parent.yml": {
			Linters: Linters{
				EnableAll: true,
				Fast:      true,
				Enable:    []string{"deadcode", "errcheck", "gosimple"},
			},
		},
		"preset_single.yml": {
			Linters: Linters{
				Fast:    true,
				Enable:  []string{"errcheck", "gosimple", "govet"},
				Disable: []string{"deadcode"},
			},
			Presets: []string{"./preset_parent.yml"},
		},
		"preset_double.yml": {
			Linters: Linters{
				Enable:  []string{"errcheck", "gosimple", "govet"},
				Disable: []string{"deadcode"},
			},
			Output: Output{
				PrintIssuedLine: true,
			},
			Presets: []string{"./preset_parent.yml", "./preset_single.yml"},
		},
		"preset_order.yml": {
			Linters: Linters{
				Fast:    true,
				Disable: []string{"deadcode", "errcheck", "gosimple", "govet"},
			},
			Presets: []string{"./preset_parent.yml", "./preset_single.yml", "./disable_all.yml"},
		},
	}

	for cfgFile, expected := range cases {
		t.Run(cfgFile, func(t *testing.T) {
			cfg := readConfig(t, cfgFile)
			assert.Equal(t, expected, cfg)
		})
	}
}

func TestFileReader_Read_Module(t *testing.T) {
	cfg := readConfig(t, "preset_module.yml")
	assert.NotEmpty(t, cfg.Linters.Enable)
}

func readConfig(t *testing.T, fn string) Config {
	var cfg, cmdCfg Config
	cmdCfg.Run.Config = filepath.Join("..", "..", "test", "testdata", "configs", fn)
	log := logutils.NewMockLog()
	log.On("Infof", mock.Anything, mock.Anything).Maybe()

	r := NewFileReader(&cfg, &cmdCfg, log)
	err := r.Read()
	require.NoError(t, err)
	return cfg
}
