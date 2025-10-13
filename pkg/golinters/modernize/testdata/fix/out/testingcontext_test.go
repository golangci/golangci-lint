//golangcitest:args -Emodernize
//golangcitest:expected_exitcode 0
package testingcontext

import (
	"context"
	"testing"
)

func Test(t *testing.T) {
	ctx := t.Context()
	_ = ctx

	func() {
		ctx, cancel := context.WithCancel(context.Background()) // Nope. scope of defer is not the testing func.
		defer cancel()
		_ = ctx
	}()

	{
		ctx := t.Context()
		_ = ctx
		var t int // not in scope of the call to WithCancel
		_ = t
	}

	{
		ctx := context.Background()
		ctx, cancel := context.WithCancel(context.Background()) // Nope. ctx is redeclared.
		defer cancel()
		_ = ctx
	}

	{
		var t int
		ctx, cancel := context.WithCancel(context.Background()) // Nope. t is shadowed.
		defer cancel()
		_ = ctx
		_ = t
	}

	t.Run("subtest", func(t2 *testing.T) {
		ctx := t2.Context()
		_ = ctx
	})
}

func TestAlt(t2 *testing.T) {
	ctx := t2.Context()
	_ = ctx
}

func Testnot(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background()) // Nope. Not a test func.
	defer cancel()
	_ = ctx
}

func Benchmark(b *testing.B) {
	ctx := b.Context()
	_ = ctx

	b.Run("subtest", func(b2 *testing.B) {
		ctx := b2.Context()
		_ = ctx
	})
}

func Fuzz(f *testing.F) {
	ctx := f.Context()
	_ = ctx
}
