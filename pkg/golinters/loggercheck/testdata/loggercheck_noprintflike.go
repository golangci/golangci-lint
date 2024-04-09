//golangcitest:args -Eloggercheck
//golangcitest:config_path testdata/loggercheck_noprintflike.yml
package loggercheck

import (
	"github.com/go-logr/logr"
)

func ExampleNoPrintfLike() {
	log := logr.Discard()

	log.Info("This message is ok")
	log.Info("Should not contains printf like format specifiers: %s %d %w") // want `logging message should not use format specifier "%s"`
	log.Info("It also checks for the key value pairs", "key", "value %.2f") // want `logging message should not use format specifier "%\.2f"`
}
