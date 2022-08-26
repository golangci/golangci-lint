//golangcitest:args -Elogrlint
//golangcitest:config_path configs/logrlint_check_logronly.yml
package logrlint

import (
	"github.com/go-logr/logr"
	"k8s.io/klog/v2"
)

func ExampleLogrOnly() {
	log := logr.Discard()
	log.Info("message", "key1", "value1", "key2", "value2", "key3") // want `odd number of arguments passed as key-value pairs for logging`

	klog.InfoS("message", "key1")
}
