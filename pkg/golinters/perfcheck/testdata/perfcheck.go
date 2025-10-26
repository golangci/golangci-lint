//golangcitest:args -Eperfcheck
package testdata

import "regexp"

func concatenate(items []string) string {
	var out string
	for _, item := range items {
		out += item // want "[perf_avoid_string_concat_loop]"
	}

	return out
}

func regexCount(inputs []string, expr string) int {
	var count int

	for _, in := range inputs {
		if regexp.MustCompile(expr).MatchString(in) { // want "[perf_regex_compile_once]"
			count++
		}
	}

	return count
}
