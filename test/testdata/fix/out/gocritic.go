//args: -Egocritic
//config: linters-settings.gocritic.enabled-checks=ruleguard
//config: linters-settings.gocritic.settings.ruleguard.rules=ruleguard/rangeExprCopy.go,ruleguard/strings_simplify.go
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
