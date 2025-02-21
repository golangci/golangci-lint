//golangcitest:args -Erevive
//golangcitest:config_path testdata/revive.yml
package testdata

import (
	"net/http"
	"time"

	"golang.org/x/xerrors"
)

func SampleRevive(t *time.Duration) error {
	if t == nil {
		return nil
	} else {
		return nil
	}
}

func testReviveComplexity(s string) { // want "cyclomatic: function testReviveComplexity has cyclomatic complexity 22"
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

func testErrorStrings() {
	_ = xerrors.New("Some error!") // want "error strings should not be capitalized or end with punctuation or a newline"
}
