package two

type Config struct {
	Version *string `yaml:"version,omitempty" toml:"version,omitempty"`

	Run Run `yaml:"run,omitempty" toml:"run,omitempty"`

	Output Output `yaml:"output,omitempty" toml:"output,omitempty"`

	Linters Linters `yaml:"linters,omitempty" toml:"linters,omitempty"`

	Issues   Issues   `yaml:"issues,omitempty" toml:"issues,omitempty"`
	Severity Severity `yaml:"severity,omitempty" toml:"severity,omitempty"`

	Formatters Formatters `yaml:"formatters,omitempty" toml:"formatters,omitempty"`
}
