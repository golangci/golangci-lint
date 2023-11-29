//go:build go1.21

//golangcitest:args -Esloglint
//golangcitest:config_path testdata/configs/sloglint_static_msg.yml
package testdata

import (
	"log/slog"
)

func test() {
	slog.Info("msg")

	const msg1 = "msg"
	slog.Info(msg1)

	msg2 := "msg"
	slog.Info(msg2) // want `message should be a string literal or a constant`
}
