//golangcitest:args -Efatcontext
//golangcitest:config_path testdata/fatcontext_checkfunctionliterals.yml
package testdata

import "context"

// Loop detection stays enabled: this MUST be reported.
func _() {
	ctx := context.Background()

	for i := 0; i < 10; i++ {
		ctx = context.WithValue(ctx, "key", i) // want "nested context in loop"
		_ = ctx
	}
}

// Function literal detection is disabled: this must NOT be reported.
func _() {
	ctx := context.Background()

	f := func() {
		ctx = context.WithValue(ctx, "key", "val")
		_ = ctx
	}
	f()
}
