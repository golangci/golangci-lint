package gosec

import (
	"testing"

	"github.com/securego/gosec/v2"
	"github.com/stretchr/testify/assert"

	"github.com/golangci/golangci-lint/pkg/config"
)

func Test_toGosecConfig(t *testing.T) {
	testCases := []struct {
		desc     string
		settings *config.GoSecSettings
		expected gosec.Config
	}{
		{
			desc:     "empty config map",
			settings: &config.GoSecSettings{},
			expected: gosec.Config{
				"global": map[gosec.GlobalOption]string{},
			},
		},
		{
			desc: "with global settings",
			settings: &config.GoSecSettings{
				Config: map[string]any{
					gosec.Globals: map[string]any{
						string(gosec.Nosec): true,
						string(gosec.Audit): "true",
					},
				},
			},
			expected: gosec.Config{
				"global": map[gosec.GlobalOption]string{
					"audit": "true",
					"nosec": "true",
				},
			},
		},
		{
			desc: "rule specified setting",
			settings: &config.GoSecSettings{
				Config: map[string]any{
					"g101": map[string]any{
						"pattern": "(?i)example",
					},
					"G301": "0750",
				},
			},
			expected: gosec.Config{
				"G101":   map[string]any{"pattern": "(?i)example"},
				"G301":   "0750",
				"global": map[gosec.GlobalOption]string{},
			},
		},
	}

	for _, test := range testCases {
		test := test
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()

			conf := toGosecConfig(test.settings)

			assert.Equal(t, test.expected, conf)
		})
	}
}
