//golangcitest:args -Edupl
//golangcitest:config_path testdata/configs/dupl.yml
package testdata

type DuplLogger struct{}

func (DuplLogger) level() int {
	return 1
}

func (DuplLogger) Debug(args ...interface{}) {}
func (DuplLogger) Info(args ...interface{})  {}

func (logger *DuplLogger) First(args ...interface{}) { // want "14-23 lines are duplicate of `.*dupl.go:25-34`"
	if logger.level() >= 0 {
		logger.Debug(args...)
		logger.Debug(args...)
		logger.Debug(args...)
		logger.Debug(args...)
		logger.Debug(args...)
		logger.Debug(args...)
	}
}

func (logger *DuplLogger) Second(args ...interface{}) { // want "25-34 lines are duplicate of `.*dupl.go:14-23`"
	if logger.level() >= 1 {
		logger.Info(args...)
		logger.Info(args...)
		logger.Info(args...)
		logger.Info(args...)
		logger.Info(args...)
		logger.Info(args...)
	}
}
