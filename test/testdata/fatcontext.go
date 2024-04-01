//golangcitest:args -Efatcontext
package testdata

import "context"

func example() {
	ctx := context.Background()

	for i := 0; i < 10; i++ {
		ctx := context.WithValue(ctx, "key", i)
		ctx = context.WithValue(ctx, "other", "val")
	}

	for i := 0; i < 10; i++ {
		ctx = context.WithValue(ctx, "key", i) // want "nested context in loop"
		ctx = context.WithValue(ctx, "other", "val")
	}

	for item := range []string{"one", "two", "three"} {
		ctx = wrapContext(ctx) // want "nested context in loop"
		ctx := context.WithValue(ctx, "key", item)
		ctx = wrapContext(ctx)
	}

	for {
		ctx = wrapContext(ctx) // want "nested context in loop"
		break
	}
}

func wrapContext(ctx context.Context) context.Context {
	return context.WithoutCancel(ctx)
}
