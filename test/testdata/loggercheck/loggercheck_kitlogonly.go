//golangcitest:args -Eloggercheck
//golangcitest:config_path configs/loggercheck_kitlogonly.yml
package loggercheck

import (
	kitlog "github.com/go-kit/log"
	"github.com/go-logr/logr"
	"go.uber.org/zap"
	"k8s.io/klog/v2"
)

func ExampleKitLogOnly_NoLogr() {
	log := logr.Discard()
	log.Info("message", "key1", "value1", "key2", "value2", "key3")
	klog.InfoS("message", "key1")

	sugar := zap.NewExample().Sugar()
	sugar.Infow("message", "key1", "value1", "key2")
	sugar.Errorw("error message", "key1")
}

func ExampleKitLogOnly() {
	logger := kitlog.NewNopLogger()

	logger.Log("msg", "message", "key1", "value1")
	logger.Log("msg")                    // want `odd number of arguments passed as key-value pairs for logging`
	logger.Log("msg", "message", "key1") // want `odd number of arguments passed as key-value pairs for logging`

	kitlog.With(logger, "key1", "value1").Log("msg", "message")
	kitlog.With(logger, "key1").Log("msg", "message") // want `odd number of arguments passed as key-value pairs for logging`
}
