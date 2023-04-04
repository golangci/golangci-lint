//golangcitest:args -Elogrlint
package loggercheck

import (
	"fmt"

	"github.com/go-logr/logr"
	"go.uber.org/zap"
	"k8s.io/klog/v2"
)

func ExampleLogrlintLogr() {
	log := logr.Discard()
	log = log.WithValues("key")                                         // want `odd number of arguments passed as key-value pairs for logging`
	log.Info("message", "key1", "value1", "key2", "value2", "key3")     // want `odd number of arguments passed as key-value pairs for logging`
	log.Error(fmt.Errorf("error"), "message", "key1", "value1", "key2") // want `odd number of arguments passed as key-value pairs for logging`
	log.Error(fmt.Errorf("error"), "message", "key1", "value1", "key2", "value2")
}

func ExampleLogrlintKlog() {
	klog.InfoS("message", "key1") // want `odd number of arguments passed as key-value pairs for logging`
}

func ExampleLogrlintZapSugar() {
	sugar := zap.NewExample().Sugar()
	defer sugar.Sync()
	sugar.Infow("message", "key1", "value1", "key2") // want `odd number of arguments passed as key-value pairs for logging`
}
