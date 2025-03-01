package two

type BaseRule struct {
	Linters    []string `yaml:"linters,omitempty" toml:"linters,omitempty"`
	Path       *string  `yaml:"path,omitempty" toml:"path,omitempty"`
	PathExcept *string  `yaml:"path-except,omitempty" toml:"path-except,omitempty"`
	Text       *string  `yaml:"text,omitempty" toml:"text,omitempty"`
	Source     *string  `yaml:"source,omitempty" toml:"source,omitempty"`

	InternalReference *string `yaml:"-,omitempty" toml:"-,omitempty"`
}
