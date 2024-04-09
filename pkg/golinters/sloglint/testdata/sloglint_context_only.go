//go:build go1.21

//golangcitest:args -Esloglint
//golangcitest:config_path testdata/sloglint_context_only.yml
package testdata

import (
	"context"
	"log/slog"
)

func test() {
	slog.InfoContext(context.Background(), "msg")

	slog.Info("msg") // want `methods without a context should not be used`
}
