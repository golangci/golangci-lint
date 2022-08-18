//golangcitest:args -Egoprintffuncname
package testdata

func PrintfLikeFuncWithBadName(format string, args ...interface{}) { // ERROR "printf-like formatting function 'PrintfLikeFuncWithBadName' should be named 'PrintfLikeFuncWithBadNamef'"
}
