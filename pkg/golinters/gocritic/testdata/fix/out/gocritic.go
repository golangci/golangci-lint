//golangcitest:args -Egocritic
//golangcitest:config_path testdata/gocritic-fix.yml
//golangcitest:expected_exitcode 0
package p

import (
	"strings"
)

func gocritic() {
	var xs [2048]byte

	// xs -> &xs
	for _, x := range &xs {
		print(x)
	}

	// strings.Count("foo", "bar") == 0 -> !strings.Contains("foo", "bar")
	if !strings.Contains("foo", "bar") {
		print("qu")
	}
}
