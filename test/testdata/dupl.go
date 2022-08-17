//golangcitest:args -Edupl
//golangcitest:config linters-settings.dupl.threshold=20
package testdata

type DuplLogger struct{}

func (DuplLogger) level() int {
	return 1
}

func (DuplLogger) Debug(args ...interface{}) {}
func (DuplLogger) Info(args ...interface{})  {}

func (logger *DuplLogger) First(args ...interface{}) { // ERROR "14-23 lines are duplicate of `.*dupl.go:25-34`"
	if logger.level() >= 0 {
		logger.Debug(args...)
		logger.Debug(args...)
		logger.Debug(args...)
		logger.Debug(args...)
		logger.Debug(args...)
		logger.Debug(args...)
	}
}

func (logger *DuplLogger) Second(args ...interface{}) { // ERROR "25-34 lines are duplicate of `.*dupl.go:14-23`"
	if logger.level() >= 1 {
		logger.Info(args...)
		logger.Info(args...)
		logger.Info(args...)
		logger.Info(args...)
		logger.Info(args...)
		logger.Info(args...)
	}
}
