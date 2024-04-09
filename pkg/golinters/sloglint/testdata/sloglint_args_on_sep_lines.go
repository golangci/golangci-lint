//go:build go1.21

//golangcitest:args -Esloglint
//golangcitest:config_path testdata/sloglint_args_on_sep_lines.yml
package testdata

import "log/slog"

func test() {
	slog.Info("msg", "foo", 1)
	slog.Info("msg",
		"foo", 1,
		"bar", 2,
	)

	slog.Info("msg", "foo", 1, "bar", 2) // want `arguments should be put on separate lines`
}
