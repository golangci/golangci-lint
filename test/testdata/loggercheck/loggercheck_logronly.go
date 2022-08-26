//golangcitest:args -Eloggercheck
//golangcitest:config_path configs/loggercheck_logronly.yml
package loggercheck

import (
	"github.com/go-logr/logr"
	"k8s.io/klog/v2"
)

func ExampleLogrOnly() {
	log := logr.Discard()
	log.Info("message", "key1", "value1", "key2", "value2", "key3") // want `odd number of arguments passed as key-value pairs for logging`

	klog.InfoS("message", "key1")
}
