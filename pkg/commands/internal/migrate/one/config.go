package one

type Config struct {
	cfgDir string // Path to the directory containing golangci-lint config file.

	Version string `mapstructure:"version"` // From v2, to be able to detect already migrated config file.

	Run Run `mapstructure:"run"`

	Output Output `mapstructure:"output"`

	LintersSettings LintersSettings `mapstructure:"linters-settings"`
	Linters         Linters         `mapstructure:"linters"`
	Issues          Issues          `mapstructure:"issues"`
	Severity        Severity        `mapstructure:"severity"`
}

func NewConfig() *Config {
	return &Config{}
}

// SetConfigDir sets the path to directory that contains golangci-lint config file.
func (c *Config) SetConfigDir(dir string) {
	c.cfgDir = dir
}

func (*Config) IsInternalTest() bool {
	return false
}
