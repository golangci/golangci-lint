//golangcitest:args -Eloggercheck
//golangcitest:config_path configs/loggercheck_custom.yml
package loggercheck

import (
	"errors"

	"go.uber.org/zap"
)

var l = New()

type Logger struct {
	s *zap.SugaredLogger
}

func New() *Logger {
	logger := zap.NewExample().Sugar()
	return &Logger{s: logger}
}

func (l *Logger) With(keysAndValues ...interface{}) *Logger {
	return &Logger{
		s: l.s.With(keysAndValues...),
	}
}

func (l *Logger) Debugw(msg string, keysAndValues ...interface{}) {
	l.s.Debugw(msg, keysAndValues...)
}

func (l *Logger) Infow(msg string, keysAndValues ...interface{}) {
	l.s.Infow(msg, keysAndValues...)
}

func (l *Logger) Warnw(msg string, keysAndValues ...interface{}) {
	l.s.Warnw(msg, keysAndValues...)
}

func (l *Logger) Errorw(msg string, keysAndValues ...interface{}) {
	l.s.Errorw(msg, keysAndValues...)
}

func (l *Logger) Fatalw(msg string, keysAndValues ...interface{}) {
	l.s.Fatalw(msg, keysAndValues...)
}

func (l *Logger) Sync() error {
	return l.s.Sync()
}

// package level wrap func

func With(keysAndValues ...interface{}) *Logger {
	return &Logger{
		s: l.s.With(keysAndValues...),
	}
}

func Debugw(msg string, keysAndValues ...interface{}) {
	l.s.Debugw(msg, keysAndValues...)
}

func Infow(msg string, keysAndValues ...interface{}) {
	l.s.Infow(msg, keysAndValues...)
}

func Warnw(msg string, keysAndValues ...interface{}) {
	l.s.Warnw(msg, keysAndValues...)
}

func Errorw(msg string, keysAndValues ...interface{}) {
	l.s.Errorw(msg, keysAndValues...)
}

func Fatalw(msg string, keysAndValues ...interface{}) {
	l.s.Fatalw(msg, keysAndValues...)
}

func Sync() error {
	return l.s.Sync()
}

func ExampleCustomLogger() {
	err := errors.New("example error")

	// custom SugaredLogger
	log := New()
	defer log.Sync()

	log.Infow("abc", "key1", "value1")
	log.Infow("abc", "key1", "value1", "key2") // want `odd number of arguments passed as key-value pairs for logging`

	log.Errorw("message", "err", err, "key1", "value1")
	log.Errorw("message", err, "key1", "value1", "key2", "value2") // want `odd number of arguments passed as key-value pairs for logging`

	// with test
	log.With("with_key1", "with_value1").Infow("message", "key1", "value1")
	log.With("with_key1", "with_value1").Infow("message", "key1", "value1", "key2") // want `odd number of arguments passed as key-value pairs for logging`
	log.With("with_key1").Infow("message", "key1", "value1")                        // want `odd number of arguments passed as key-value pairs for logging`
}

func ExampleCustomLoggerPackageLevelFunc() {
	err := errors.New("example error")

	defer Sync()

	Infow("abc", "key1", "value1")
	Infow("abc", "key1", "value1", "key2") // want `odd number of arguments passed as key-value pairs for logging`

	Errorw("message", "err", err, "key1", "value1")
	Errorw("message", err, "key1", "value1", "key2", "value2") // want `odd number of arguments passed as key-value pairs for logging`

	// with test
	With("with_key1", "with_value1").Infow("message", "key1", "value1")
	With("with_key1", "with_value1").Infow("message", "key1", "value1", "key2") // want `odd number of arguments passed as key-value pairs for logging`
	With("with_key1").Infow("message", "key1", "value1")                        // want `odd number of arguments passed as key-value pairs for logging`
}
