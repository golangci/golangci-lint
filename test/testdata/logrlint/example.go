//golangcitest:args -Elogrlint
package logrlint

import (
	"fmt"

	"github.com/go-logr/logr"
)

func Example() {
	log := logr.Discard()
	log = log.WithValues("key")                                         // Error `odd number of arguments passed as key-value pairs for logging`
	log.Info("message", "key1", "value1", "key2", "value2", "key3")     // Error `odd number of arguments passed as key-value pairs for logging`
	log.Error(fmt.Errorf("error"), "message", "key1", "value1", "key2") // Error `odd number of arguments passed as key-value pairs for logging`
	log.Error(fmt.Errorf("error"), "message", "key1", "value1", "key2", "value2")
}
