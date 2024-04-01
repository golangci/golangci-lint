//golangcitest:args -Egoprintffuncname
package testdata

func PrintfLikeFuncWithBadName(format string, args ...interface{}) { // want "printf-like formatting function 'PrintfLikeFuncWithBadName' should be named 'PrintfLikeFuncWithBadNamef'"
}
