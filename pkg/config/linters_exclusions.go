package config

type LinterExclusions struct {
	Generated  string        `mapstructure:"generated"`
	WarnUnused bool          `mapstructure:"warn-unused"`
	Default    []string      `mapstructure:"default"`
	Rules      []ExcludeRule `mapstructure:"rules"`
	Paths      []string      `mapstructure:"paths"`
}
