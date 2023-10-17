//go:build go1.21

//golangcitest:args -Esloglint
//golangcitest:config_path testdata/configs/sloglint_no_raw_keys.yml
package testdata

import "log/slog"

const foo = "foo"

func Foo(value int) slog.Attr {
	return slog.Int("foo", value)
}

func test() {
	slog.Info("msg", foo, 1)
	slog.Info("msg", Foo(1))

	slog.Info("msg", "foo", 1)           // want `raw keys should not be used`
	slog.Info("msg", slog.Int("foo", 1)) // want `raw keys should not be used`
}
