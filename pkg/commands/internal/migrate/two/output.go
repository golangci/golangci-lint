package two

type Output struct {
	Formats    Formats  `yaml:"formats,omitempty" toml:"formats,omitempty"`
	SortOrder  []string `yaml:"sort-order,omitempty" toml:"sort-order,omitempty"`
	PathPrefix *string  `yaml:"path-prefix,omitempty" toml:"path-prefix,omitempty"`
	ShowStats  *bool    `yaml:"show-stats,omitempty" toml:"show-stats,omitempty"`
}
