TODO(ldez): update doc about plugin inside `.golangci.reference.yml`

```go
// pkg/lint/lintersdb/builder_plugin.go

const goPluginType = "goplugin"

// Build loads private linters that are specified in the golangci config file.
func (b *PluginBuilder) Build(cfg *config.Config) []*linter.Config {
	if cfg == nil || b.log == nil {
		return nil
	}

	var linters []*linter.Config

	for name, settings := range cfg.LintersSettings.Custom {
		if settings.Type != goPluginType && settings.Type != "" {
			continue
		}

		lc, err := b.loadConfig(cfg, name, settings)
		if err != nil {
			b.log.Errorf("Unable to load custom analyzer %s:%s, %v", name, settings.Path, err)
		} else {
			linters = append(linters, lc)
		}
	}

	return linters
}
```

```go
// pkg/config/linters_settings.go

type CustomLinterSettings struct {
	// Path to a plugin *.so file that implements the private linter.
	Path string
	// Description describes the purpose of the private linter.
	Description string
	// OriginalURL The URL containing the source code for the private linter.
	OriginalURL string `mapstructure:"original-url"`

	// Settings plugin settings only work with linterdb.PluginConstructor symbol.
	Settings any

	// FIXME goplugin,module
	Type string `mapstructure:"type"`
}

func (s *CustomLinterSettings) Validate() error {
	if s.Type == "module" {
		return nil
	}

	if s.Path == "" {
		return errors.New("path is required")
	}

	return nil
}

```

```go
// pkg/config/linters_settings.go

func (s *LintersSettings) Validate() error {
	if err := s.Govet.Validate(); err != nil {
		return err
	}

	for name, settings := range s.Custom {
		if err := settings.Validate(); err != nil {
			return fmt.Errorf("custom linter %q: %w", name, err)
		}
	}

	return nil
}

```
