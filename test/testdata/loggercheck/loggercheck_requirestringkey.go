//golangcitest:args -Eloggercheck
//golangcitest:config_path configs/loggercheck_requirestringkey.yml
package loggercheck

import (
	"github.com/go-logr/logr"
)

func ExampleRequireStringKey() {
	log := logr.Discard()
	log.Info("message", "key1", "value1")
	const key1 = "key1"
	log.Info("message", key1, "value1")

	key2 := []byte(key1)
	log.Info("message", key2, "value2") // want `logging keys are expected to be inlined constant strings, please replace "key2" provided with string`

	key3 := key1
	log.Info("message", key3, "value3") // want `logging keys are expected to be inlined constant strings, please replace "key3" provided with string`

	log.Info("message", "键1", "value1") // want `logging keys are expected to be alphanumeric strings, please remove any non-latin characters from "键1"`
}
