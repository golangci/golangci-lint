//go:build go1.21

//golangcitest:args -Esloglint
package testdata

import "log/slog"

func test() {
	slog.Info("msg", "foo", 1, "bar", 2)
	slog.Info("msg", slog.Int("foo", 1), slog.Int("bar", 2))

	slog.Info("msg", "foo", 1, slog.Int("bar", 2)) // want `key-value pairs and attributes should not be mixed`
}
