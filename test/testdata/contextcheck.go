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

	funcWithCtx(context.Background()) // ERROR `The context param may be context.TODO\(\) or context.Background\(\), please replace it with another way, such as context.WithValue\(ctx, key, val\)`
}

func contextcheckCase3(ctx context.Context) {
	func() {
		funcWithCtx(ctx)
	}()

	ctx = context.Background() // ERROR `Invalid call to get new context, please replace it with another way, such as context.WithValue\(ctx, key, val\)`
	funcWithCtx(ctx)
}

func funcWithCtx(ctx context.Context) {}

func funcWithoutCtx() {
	funcWithCtx(context.TODO())
}
