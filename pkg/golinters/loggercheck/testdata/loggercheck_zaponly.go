//golangcitest:args -Eloggercheck
//golangcitest:config_path testdata/loggercheck_zaponly.yml
package loggercheck

import "go.uber.org/zap"

func ExampleZapOnly() {
	sugar := zap.NewExample().Sugar()

	sugar.Infow("message", "key1", "value1", "key2") // want `odd number of arguments passed as key-value pairs for logging`
	sugar.Errorw("error message", "key1")            // want `odd number of arguments passed as key-value pairs for logging`
}
