//args: -Econtextcheck
package testdata

import "context"

type MyString string

func contextcheckCase1(ctx context.Context) {
	funcWithoutCtx() // ERROR "Function `funcWithoutCtx` should pass the context parameter"
}

func contextcheckCase2(ctx context.Context) {
	ctx = context.WithValue(ctx, MyString("aaa"), "aaaaaa")
	funcWithCtx(ctx)

	defer func() {
		funcWithCtx(ctx)
	}()

	func(ctx context.Context) {
		funcWithCtx(ctx)
	}(ctx)

	funcWithCtx(context.Background()) // ERROR "Non-inherited new context, use function like `context.WithXXX` instead"
}

func contextcheckCase3(ctx context.Context) {
	func() {
		funcWithCtx(ctx)
	}()

	ctx = context.Background() // ERROR "Non-inherited new context, use function like `context.WithXXX` instead"
	funcWithCtx(ctx)
}

func contextcheckCase4(ctx context.Context) {
	ctx, cancel := getNewCtx(ctx)
	defer cancel()
	funcWithCtx(ctx)
}

func funcWithCtx(ctx context.Context) {}

func funcWithoutCtx() {
	funcWithCtx(context.TODO())
}

func getNewCtx(ctx context.Context) (newCtx context.Context, cancel context.CancelFunc) {
	return context.WithCancel(ctx)
}
