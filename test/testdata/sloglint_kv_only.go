//go:build go1.21

//golangcitest:args -Esloglint
//golangcitest:config_path testdata/configs/sloglint_kv_only.yml
package testdata

import "log/slog"

func test() {
	slog.Info("msg", "foo", 1, "bar", 2)

	slog.Info("msg", slog.Int("foo", 1), slog.Int("bar", 2)) // want `attributes should not be used`
}
