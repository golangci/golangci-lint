package internal

import (
	"fmt"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type FlagFunc[T any] func(name string, value T, usage string) *T

type FlagPFunc[T any] func(name, shorthand string, value T, usage string) *T

// AddFlagAndBind adds a Cobra/pflag flag and binds it with Viper.
func AddFlagAndBind[T any](v *viper.Viper, fs *pflag.FlagSet, pfn FlagFunc[T], name, bind string, value T, usage string) {
	pfn(name, value, usage)

	err := v.BindPFlag(bind, fs.Lookup(name))
	if err != nil {
		panic(fmt.Sprintf("failed to bind flag %s: %v", name, err))
	}
}

// AddFlagAndBindP adds a Cobra/pflag flag and binds it with Viper.
func AddFlagAndBindP[T any](v *viper.Viper, fs *pflag.FlagSet, pfn FlagPFunc[T], name, shorthand, bind string, value T, usage string) {
	pfn(name, shorthand, value, usage)

	err := v.BindPFlag(bind, fs.Lookup(name))
	if err != nil {
		panic(fmt.Sprintf("failed to bind flag %s: %v", name, err))
	}
}
