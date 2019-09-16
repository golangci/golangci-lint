package log

// Interface represents the API of both Logger and Entry.
type Interface interface {
	WithFields(fields Fielder) *Entry
	WithField(key string, value interface{}) *Entry
	WithError(err error) *Entry
	Debug(msg string)
	Info(msg string)
	Warn(msg string)
	Error(msg string)
	Fatal(msg string)
	Debugf(msg string, v ...interface{})
	Infof(msg string, v ...interface{})
	Warnf(msg string, v ...interface{})
	Errorf(msg string, v ...interface{})
	Fatalf(msg string, v ...interface{})
	Trace(msg string) *Entry
}
