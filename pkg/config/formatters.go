package config

type Formatters struct {
	Enable     []string            `mapstructure:"enable"`
	Settings   FormatterSettings   `mapstructure:"settings"`
	Exclusions FormatterExclusions `mapstructure:"exclusions"`
}

type FormatterExclusions struct {
	Generated string   `mapstructure:"generated"`
	Paths     []string `mapstructure:"paths"`
}
