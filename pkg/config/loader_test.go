package config

import (
	"testing"

	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/golangci/golangci-lint/v2/pkg/logutils"
)

func TestLoader_handleClearConfigOutputs(t *testing.T) {
	t.Run("flag not set", func(t *testing.T) {
		// Setup
		cfg := &Config{
			Output: Output{
				Formats: Formats{
					JSON: SimpleFormat{Path: "/tmp/config.json"},
					HTML: SimpleFormat{Path: "/tmp/config.html"},
				},
			},
		}

		fs := pflag.NewFlagSet("test", pflag.ContinueOnError)
		fs.Bool("clear-config-outputs", false, "test flag")

		loader := &Loader{
			BaseLoader: &BaseLoader{
				log: logutils.NewStderrLog(logutils.DebugKeyEmpty),
			},
			fs:  fs,
			cfg: cfg,
		}

		// Execute
		err := loader.handleClearConfigOutputs()
		require.NoError(t, err)

		// Verify - config outputs should remain unchanged
		assert.Equal(t, "/tmp/config.json", cfg.Output.Formats.JSON.Path)
		assert.Equal(t, "/tmp/config.html", cfg.Output.Formats.HTML.Path)
	})

	t.Run("flag set with no CLI outputs", func(t *testing.T) {
		// Setup
		cfg := &Config{
			Output: Output{
				Formats: Formats{
					JSON: SimpleFormat{Path: "/tmp/config.json"},
					HTML: SimpleFormat{Path: "/tmp/config.html"},
				},
			},
		}

		fs := pflag.NewFlagSet("test", pflag.ContinueOnError)
		fs.Bool("clear-config-outputs", false, "test flag")
		fs.String("output.json.path", "", "json output path")
		fs.String("output.html.path", "", "html output path")

		// Set the flag
		err := fs.Set("clear-config-outputs", "true")
		require.NoError(t, err)

		loader := &Loader{
			BaseLoader: &BaseLoader{
				log: logutils.NewStderrLog(logutils.DebugKeyEmpty),
			},
			fs:  fs,
			cfg: cfg,
		}

		// Execute
		err = loader.handleClearConfigOutputs()
		require.NoError(t, err)

		// Verify - all config outputs should be cleared
		assert.Empty(t, cfg.Output.Formats.JSON.Path)
		assert.Empty(t, cfg.Output.Formats.HTML.Path)
	})

	t.Run("flag set with CLI JSON output", func(t *testing.T) {
		// Setup
		cfg := &Config{
			Output: Output{
				Formats: Formats{
					JSON: SimpleFormat{Path: "/tmp/config.json"},
					HTML: SimpleFormat{Path: "/tmp/config.html"},
				},
			},
		}

		fs := pflag.NewFlagSet("test", pflag.ContinueOnError)
		fs.Bool("clear-config-outputs", false, "test flag")
		fs.String("output.json.path", "", "json output path")
		fs.String("output.html.path", "", "html output path")

		// Set the flag and CLI output
		err := fs.Set("clear-config-outputs", "true")
		require.NoError(t, err)
		err = fs.Set("output.json.path", "/tmp/cli.json")
		require.NoError(t, err)

		loader := &Loader{
			BaseLoader: &BaseLoader{
				log: logutils.NewStderrLog(logutils.DebugKeyEmpty),
			},
			fs:  fs,
			cfg: cfg,
		}

		// Execute
		err = loader.handleClearConfigOutputs()
		require.NoError(t, err)

		// Verify - only CLI output should remain
		assert.Equal(t, "/tmp/cli.json", cfg.Output.Formats.JSON.Path)
		assert.Empty(t, cfg.Output.Formats.HTML.Path)
	})

	t.Run("flag set with multiple CLI outputs", func(t *testing.T) {
		// Setup
		cfg := &Config{
			Output: Output{
				Formats: Formats{
					JSON: SimpleFormat{Path: "/tmp/config.json"},
					HTML: SimpleFormat{Path: "/tmp/config.html"},
					Text: Text{
						SimpleFormat: SimpleFormat{Path: "/tmp/config.txt"},
					},
				},
			},
		}

		fs := pflag.NewFlagSet("test", pflag.ContinueOnError)
		fs.Bool("clear-config-outputs", false, "test flag")
		fs.String("output.json.path", "", "json output path")
		fs.String("output.html.path", "", "html output path")
		fs.String("output.text.path", "", "text output path")

		// Set the flag and CLI outputs
		err := fs.Set("clear-config-outputs", "true")
		require.NoError(t, err)
		err = fs.Set("output.json.path", "/tmp/cli.json")
		require.NoError(t, err)
		err = fs.Set("output.html.path", "/tmp/cli.html")
		require.NoError(t, err)

		loader := &Loader{
			BaseLoader: &BaseLoader{
				log: logutils.NewStderrLog(logutils.DebugKeyEmpty),
			},
			fs:  fs,
			cfg: cfg,
		}

		// Execute
		err = loader.handleClearConfigOutputs()
		require.NoError(t, err)

		// Verify - only CLI outputs should remain
		assert.Equal(t, "/tmp/cli.json", cfg.Output.Formats.JSON.Path)
		assert.Equal(t, "/tmp/cli.html", cfg.Output.Formats.HTML.Path)
		assert.Empty(t, cfg.Output.Formats.Text.Path)
	})

	t.Run("flag set with CLI format options", func(t *testing.T) {
		// Setup
		cfg := &Config{
			Output: Output{
				Formats: Formats{
					Text: Text{
						SimpleFormat:    SimpleFormat{Path: "/tmp/config.txt"},
						PrintLinterName: false,
						PrintIssuedLine: false,
						Colors:          false,
					},
				},
			},
		}

		fs := pflag.NewFlagSet("test", pflag.ContinueOnError)
		fs.Bool("clear-config-outputs", false, "test flag")
		fs.String("output.text.path", "", "text output path")
		fs.Bool("output.text.print-linter-name", true, "print linter name")
		fs.Bool("output.text.colors", true, "use colors")

		// Set the flag and CLI outputs with options
		err := fs.Set("clear-config-outputs", "true")
		require.NoError(t, err)
		err = fs.Set("output.text.path", "/tmp/cli.txt")
		require.NoError(t, err)
		err = fs.Set("output.text.print-linter-name", "true")
		require.NoError(t, err)
		err = fs.Set("output.text.colors", "false")
		require.NoError(t, err)

		loader := &Loader{
			BaseLoader: &BaseLoader{
				log: logutils.NewStderrLog(logutils.DebugKeyEmpty),
			},
			fs:  fs,
			cfg: cfg,
		}

		// Execute
		err = loader.handleClearConfigOutputs()
		require.NoError(t, err)

		// Verify - CLI output with options should be preserved
		assert.Equal(t, "/tmp/cli.txt", cfg.Output.Formats.Text.Path)
		assert.True(t, cfg.Output.Formats.Text.PrintLinterName)
		assert.False(t, cfg.Output.Formats.Text.Colors)
		assert.False(t, cfg.Output.Formats.Text.PrintIssuedLine) // Not set via CLI, should be default false
	})
}
