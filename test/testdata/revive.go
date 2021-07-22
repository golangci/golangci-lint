//args: -Erevive
//config_path: testdata/configs/revive.yml
package testdata

import (
	"net/http"
	"time"
)

func SampleRevive(t *time.Duration) error {
	if t == nil {
		return nil
	} else {
		return nil
	}
}

func testReviveComplexity(s string) { // ERROR "cyclomatic: function testReviveComplexity has cyclomatic complexity 22"
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
