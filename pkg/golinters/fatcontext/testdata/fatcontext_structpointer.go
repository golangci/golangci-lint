//golangcitest:args -Efatcontext
//golangcitest:config_path testdata/fatcontext_structpointer.yml
package testdata

import (
	"context"
)

type Container struct {
	Ctx context.Context
}

func something() func(*Container) {
	return func(r *Container) {
		ctx := r.Ctx
		ctx = context.WithValue(ctx, "key", "val")
		r.Ctx = ctx // want "potential nested context in struct pointer"
	}
}

func blah(r *Container) {
	ctx := r.Ctx
	ctx = context.WithValue(ctx, "key", "val")
	r.Ctx = ctx // want "potential nested context in struct pointer"
}
