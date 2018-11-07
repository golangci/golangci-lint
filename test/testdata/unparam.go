//args: -Eunparam
package testdata

func unparamUnused(a, b uint) uint { // ERROR "`unparamUnused` - `b` is unused"
	a++
	return a
}
