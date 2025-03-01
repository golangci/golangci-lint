package two

type LinterExclusions struct {
	Generated   *string       `yaml:"generated,omitempty" toml:"generated,omitempty"`
	WarnUnused  *bool         `yaml:"warn-unused,omitempty" toml:"warn-unused,omitempty"`
	Presets     []string      `yaml:"presets,omitempty" toml:"presets,omitempty"`
	Rules       []ExcludeRule `yaml:"rules,omitempty" toml:"rules,omitempty"`
	Paths       []string      `yaml:"paths,omitempty" toml:"paths,omitempty"`
	PathsExcept []string      `yaml:"paths-except,omitempty" toml:"paths-except,omitempty"`
}

type ExcludeRule struct {
	BaseRule `yaml:",inline"`
}
