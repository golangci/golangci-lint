//args: -Eargslen
//config: linters-settings.argslen.maxArguments=1
package testdata

func argslenFunc(s1, s2, s3 string) { // ERROR "args number for function `argslenFunc` is too high, 3 vs 1 allowed"
	return
}
