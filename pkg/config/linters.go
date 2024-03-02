package config

import (
	"errors"
	"fmt"
)

type Linters struct {
	Enable     []string
	Disable    []string
	EnableAll  bool `mapstructure:"enable-all"`
	DisableAll bool `mapstructure:"disable-all"`
	Fast       bool

	Presets []string
}

func (l *Linters) Validate() error {
	if err := l.validateAllDisableEnableOptions(); err != nil {
		return err
	}

	if err := l.validateDisabledAndEnabledAtOneMoment(); err != nil {
		return err
	}

	return nil
}

func (l *Linters) validateAllDisableEnableOptions() error {
	if l.EnableAll && l.DisableAll {
		return errors.New("--enable-all and --disable-all options must not be combined")
	}

	if l.DisableAll {
		if len(l.Enable) == 0 && len(l.Presets) == 0 {
			return errors.New("all linters were disabled, but no one linter was enabled: must enable at least one")
		}

		if len(l.Disable) != 0 {
			return errors.New("can't combine options --disable-all and --disable")
		}
	}

	if l.EnableAll && len(l.Enable) != 0 && !l.Fast {
		return errors.New("can't combine options --enable-all and --enable")
	}

	return nil
}

func (l *Linters) validateDisabledAndEnabledAtOneMoment() error {
	enabledLintersSet := map[string]bool{}
	for _, name := range l.Enable {
		enabledLintersSet[name] = true
	}

	for _, name := range l.Disable {
		if enabledLintersSet[name] {
			return fmt.Errorf("linter %q can't be disabled and enabled at one moment", name)
		}
	}

	return nil
}
