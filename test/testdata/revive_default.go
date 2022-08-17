//golangcitest:args -Erevive
package testdata

import (
	"net/http"
	"time"
)

func testReviveDefault(t *time.Duration) error {
	if t == nil {
		return nil
	} else { // ERROR "indent-error-flow: if block ends with a return statement, .*"
		return nil
	}
}

func testReviveComplexityDefault(s string) {
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
