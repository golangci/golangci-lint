//args: -Edupl --dupl.threshold=20
package testdata

type DuplLogger struct{}

func (DuplLogger) level() int {
	return 1
}

func (DuplLogger) Debug(args ...interface{}) {}
func (DuplLogger) Info(args ...interface{})  {}

func (logger *DuplLogger) First(args ...interface{}) { // ERROR "13-22 lines are duplicate of `testdata/dupl.go:24-33`"
	if logger.level() >= 0 {
		logger.Debug(args...)
		logger.Debug(args...)
		logger.Debug(args...)
		logger.Debug(args...)
		logger.Debug(args...)
		logger.Debug(args...)
	}
}

func (logger *DuplLogger) Second(args ...interface{}) { // ERROR "24-33 lines are duplicate of `testdata/dupl.go:13-22`"
	if logger.level() >= 1 {
		logger.Info(args...)
		logger.Info(args...)
		logger.Info(args...)
		logger.Info(args...)
		logger.Info(args...)
		logger.Info(args...)
	}
}
