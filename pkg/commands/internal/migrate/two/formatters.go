package two

type Formatters struct {
	Enable     []string            `yaml:"enable,omitempty" toml:"enable,omitempty"`
	Settings   FormatterSettings   `yaml:"settings,omitempty" toml:"settings,omitempty"`
	Exclusions FormatterExclusions `yaml:"exclusions,omitempty" toml:"exclusions,omitempty"`
}

type FormatterExclusions struct {
	Generated *string  `yaml:"generated,omitempty" toml:"generated,omitempty"`
	Paths     []string `yaml:"paths,omitempty" toml:"paths,omitempty"`
}
