//golangcitest:args -Einamedparam
//golangcitest:config_path testdata/inamedparam_skip_single_param.yml
package testdata

import "context"

type NamedParam interface {
	Void()

	SingleParam(string) error

	WithName(ctx context.Context, number int, toggle bool) (bool, error)

	WithoutName(
		context.Context, // want "interface method WithoutName must have named param for type context.Context"
		int, // want "interface method WithoutName must have named param for type int"
	)
}
