package golinters

import (
	"fmt"
	"testing"

	"github.com/securego/gosec/v2"
	"github.com/stretchr/testify/assert"

	"github.com/golangci/golangci-lint/pkg/config"
)

func TestToGosecConfig(t *testing.T) {
	t.Run("empty config map", func(t *testing.T) {
		settings := &config.GoSecSettings{}

		gosecConfig := toGosecConfig(settings)
		assert.Len(t, gosecConfig, 1)
		assert.Contains(t, gosecConfig, gosec.Globals)
	})

	t.Run("with global settings", func(t *testing.T) {
		globalsSettings := map[string]any{
			string(gosec.Nosec): true,
			string(gosec.Audit): "true",
		}
		settings := &config.GoSecSettings{
			Config: map[string]any{
				gosec.Globals: globalsSettings,
			},
		}

		gosecConfig := toGosecConfig(settings)
		assert.Len(t, gosecConfig, 1)
		assert.Contains(t, gosecConfig, gosec.Globals)

		for _, k := range []gosec.GlobalOption{gosec.Nosec, gosec.Audit} {
			v, err := gosecConfig.GetGlobal(k)
			assert.NoError(t, err, "error getting global option %s", k)
			assert.Equal(
				t,
				fmt.Sprintf("%v", globalsSettings[string(k)]),
				v,
				"global option %s is not set to expected value", k,
			)
		}

		for _, k := range []gosec.GlobalOption{gosec.NoSecAlternative} {
			_, err := gosecConfig.GetGlobal(k)
			assert.Error(t, err, "should not set global option %s", k)
		}
	})

	t.Run("rule specified settings", func(t *testing.T) {
		settings := &config.GoSecSettings{
			Config: map[string]any{
				"g101": map[string]any{
					"pattern": "(?i)example",
				},
				"G301": "0750",
			},
		}

		gosecConfig := toGosecConfig(settings)
		assert.Equal(t,
			gosec.Config{
				"G101":   map[string]any{"pattern": "(?i)example"},
				"G301":   "0750",
				"global": map[gosec.GlobalOption]string{}},
			gosecConfig)
	})
}
