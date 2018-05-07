package testdata

type DuplLogger struct{}

func (DuplLogger) level() int {
	return 1
}

func (DuplLogger) Debug(args ...interface{}) {}
func (DuplLogger) Info(args ...interface{})  {}

func (logger *DuplLogger) First(args ...interface{}) { // ERROR "12-21 lines are duplicate of `testdata/with_issues/dupl.go:23-32`"
	if logger.level() >= 0 {
		logger.Debug(args...)
		logger.Debug(args...)
		logger.Debug(args...)
		logger.Debug(args...)
		logger.Debug(args...)
		logger.Debug(args...)
	}
}

func (logger *DuplLogger) Second(args ...interface{}) { // ERROR "23-32 lines are duplicate of `testdata/with_issues/dupl.go:12-21`"
	if logger.level() >= 1 {
		logger.Info(args...)
		logger.Info(args...)
		logger.Info(args...)
		logger.Info(args...)
		logger.Info(args...)
		logger.Info(args...)
	}
}
