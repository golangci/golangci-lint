//golangcitest:args -Eloggercheck
//golangcitest:config_path testdata/loggercheck_slogonly.yml
package loggercheck

import "log/slog"

func ExampleSlogOnly() {
	logger := slog.With("key1", "value1")

	logger.Info("msg", "key1") // want `odd number of arguments passed as key-value pairs for logging`
	slog.Info("msg", "key1")   // want `odd number of arguments passed as key-value pairs for logging`
}
