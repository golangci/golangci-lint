//go:build go1.21

//golangcitest:args -Esloglint
//golangcitest:config_path testdata/sloglint_forbidden_keys_case.yml
package testdata

import "log/slog"

const (
	allowedKey   = "foo_bar"
	forbiddenKey = "foo-bar"
)

func test() {
	slog.Info("msg", "foo_bar", 1)
	slog.Info("msg", allowedKey, 1)
	slog.Info("msg", slog.Int("foo_bar", 1))
	slog.Info("msg", slog.Int(allowedKey, 1))

	slog.Info("msg", "foo-bar", 1)               // want `"foo-bar" key is forbidden and should not be used`
	slog.Info("msg", forbiddenKey, 1)            // want `"foo-bar" key is forbidden and should not be used`
	slog.Info("msg", slog.Int("foo-bar", 1))     // want `"foo-bar" key is forbidden and should not be used`
	slog.Info("msg", slog.Int(forbiddenKey, 1))  // want `"foo-bar" key is forbidden and should not be used`
	slog.Info("msg", slog.Int("foo-bar-baz", 1)) // want `"foo-bar-baz" key is forbidden and should not be used`

	slog.With("foo-bar", 1).Info("msg")               // want `"foo-bar" key is forbidden and should not be used`
	slog.With(forbiddenKey, 1).Info("msg")            // want `"foo-bar" key is forbidden and should not be used`
	slog.With(slog.Int("foo-bar", 1)).Info("msg")     // want `"foo-bar" key is forbidden and should not be used`
	slog.With(slog.Int(forbiddenKey, 1)).Info("msg")  // want `"foo-bar" key is forbidden and should not be used`
	slog.With(slog.Int("foo-bar-baz", 1)).Info("msg") // want `"foo-bar-baz" key is forbidden and should not be used`
}
