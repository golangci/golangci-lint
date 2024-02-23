package internal

import (
	"fmt"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type FlagFunc[T any] func(name string, value T, usage string) *T

type FlagPFunc[T any] func(name, shorthand string, value T, usage string) *T

// Vibra adds a Cobra flag and binds it with Viper.
func Vibra[T any](v *viper.Viper, fs *pflag.FlagSet, pfn FlagFunc[T], name, bind string, value T, usage string) {
	pfn(name, value, usage)

	err := v.BindPFlag(bind, fs.Lookup(name))
	if err != nil {
		panic(fmt.Sprintf("failed to bind flag %s: %v", name, err))
	}
}

// VibraP adds a Cobra flag and binds it with Viper.
func VibraP[T any](v *viper.Viper, fs *pflag.FlagSet, pfn FlagPFunc[T], name, shorthand, bind string, value T, usage string) {
	pfn(name, shorthand, value, usage)

	err := v.BindPFlag(bind, fs.Lookup(name))
	if err != nil {
		panic(fmt.Sprintf("failed to bind flag %s: %v", name, err))
	}
}
