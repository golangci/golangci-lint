package two

type Linters struct {
	Default  *string  `yaml:"default,omitempty" toml:"default,omitempty"`
	Enable   []string `yaml:"enable,omitempty" toml:"enable,omitempty"`
	Disable  []string `yaml:"disable,omitempty" toml:"disable,omitempty"`
	FastOnly *bool    `yaml:"fast-only,omitempty" toml:"fast-only,omitempty"`

	Settings LintersSettings `yaml:"settings,omitempty" toml:"settings,omitempty"`

	Exclusions LinterExclusions `yaml:"exclusions,omitempty" toml:"exclusions,omitempty"`
}
