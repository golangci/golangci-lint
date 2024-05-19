//golangcitest:args -Esloglint
//golangcitest:config_path testdata/sloglint_forbidden_keys.yml
package testdata

import "log/slog"

const (
	snakeKey = "foo_bar"
)

func tests() {
	slog.Info("msg")
	slog.Info("msg", "foo-bar", 1)
	slog.Info("msg", "foo_bar", 1)           // want `"foo_bar" key is forbidden and should not be used`
	slog.Info("msg", snakeKey, 1)            // want `"foo_bar" key is forbidden and should not be used`
	slog.Info("msg", slog.Int("foo_bar", 1)) // want `"foo_bar" key is forbidden and should not be used`
	slog.Info("msg", slog.Int(snakeKey, 1))  // want `"foo_bar" key is forbidden and should not be used`
	slog.Info("msg", slog.Attr{})
	slog.Info("msg", slog.Attr{"foo_bar", slog.IntValue(1)})             // want `"foo_bar" key is forbidden and should not be used`
	slog.Info("msg", slog.Attr{snakeKey, slog.IntValue(1)})              // want `"foo_bar" key is forbidden and should not be used`
	slog.Info("msg", slog.Attr{Key: "foo_bar"})                          // want `"foo_bar" key is forbidden and should not be used`
	slog.Info("msg", slog.Attr{Key: snakeKey})                           // want `"foo_bar" key is forbidden and should not be used`
	slog.Info("msg", slog.Attr{Key: "foo_bar", Value: slog.IntValue(1)}) // want `"foo_bar" key is forbidden and should not be used`
	slog.Info("msg", slog.Attr{Key: snakeKey, Value: slog.IntValue(1)})  // want `"foo_bar" key is forbidden and should not be used`
	slog.Info("msg", slog.Attr{Value: slog.IntValue(1), Key: "foo_bar"}) // want `"foo_bar" key is forbidden and should not be used`
	slog.Info("msg", slog.Attr{Value: slog.IntValue(1), Key: snakeKey})  // want `"foo_bar" key is forbidden and should not be used`
	slog.Info("msg", slog.Attr{Value: slog.IntValue(1), Key: `foo_bar`}) // want `"foo_bar" key is forbidden and should not be used`
}
