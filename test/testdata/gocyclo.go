//args: -Egocyclo
//config: linters-settings.gocyclo.min-complexity=20
package testdata

import "net/http"

func GocycloBigComplexity(s string) { // ERROR "cyclomatic complexity .* of func .* is high .*"
	if s == http.MethodGet || s == "2" || s == "3" || s == "4" || s == "5" || s == "6" || s == "7" {
		return
	}

	if s == "1" || s == "2" || s == "3" || s == "4" || s == "5" || s == "6" || s == "7" {
		return
	}

	if s == "1" || s == "2" || s == "3" || s == "4" || s == "5" || s == "6" || s == "7" {
		return
	}
}
