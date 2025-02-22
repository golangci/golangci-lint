package config

const (
	GroupStandard = "standard"
	GroupAll      = "all"
	GroupNone     = "none"
	GroupFast     = "fast"
)

type Linters struct {
	Default string   `mapstructure:"default"`
	Enable  []string `mapstructure:"enable"`
	Disable []string `mapstructure:"disable"`

	Settings LintersSettings `mapstructure:"settings"`

	Exclusions LinterExclusions `mapstructure:"exclusions"`
}

func (l *Linters) Validate() error {
	validators := []func() error{
		l.Exclusions.Validate,
	}

	for _, v := range validators {
		if err := v(); err != nil {
			return err
		}
	}

	return nil
}
