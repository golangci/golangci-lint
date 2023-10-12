//golangcitest:args -Einamedparam
package testdata

import "context"

type tStruct struct {
	a int
}

type Doer interface {
	Do() string
}

type NamedParam interface {
	Void()

	NoArgs() string

	WithName(ctx context.Context, number int, toggle bool, tStruct *tStruct, doer Doer) (bool, error)

	WithoutName(
		context.Context, // want "interface method WithoutName must have named param for type context.Context"
		int, // want "interface method WithoutName must have named param for type int"
		bool, // want "interface method WithoutName must have named param for type bool"
		tStruct, // want "interface method WithoutName must have named param for type tStruct"
		Doer, // want "interface method WithoutName must have named param for type Doer"
		struct{ b bool }, // want "interface method WithoutName must have all named params"
	)
}
