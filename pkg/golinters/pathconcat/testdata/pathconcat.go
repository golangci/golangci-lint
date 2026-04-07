//golangcitest:args -Epathconcat
package testdata

import (
	"fmt"
	"strings"
)

func concatWithSlash(a, b string) string {
	return a + "/" + b // want `use path\.Join\(\) instead of string concatenation with "/"`
}

func concatMultiSegment(a, b, c string) string {
	return a + "/" + b + "/" + c // want `use path\.Join\(\) instead of string concatenation with "/"`
}

func sprintfPath(a, b string) string {
	return fmt.Sprintf("%s/%s", a, b) // want `use path\.Join\(\) instead of fmt\.Sprintf with path separators`
}

func stringsJoinSlash(parts []string) string {
	return strings.Join(parts, "/") // want `use path\.Join\(\) instead of strings\.Join with "/"`
}

func schemePrefix(host string) string {
	return "https://" + host // OK: scheme prefix
}

func regularConcat(a, b string) string {
	return a + b // OK: no slash
}

func concatNonSep(a string) string {
	return a + "/api" // OK: no bare "/" separator
}

func sprintfNoPathSep(a string) string {
	return fmt.Sprintf("value: %s", a) // OK: no path separator
}

func stringsJoinComma(parts []string) string {
	return strings.Join(parts, ",") // OK: not a slash
}

func postgresConnStr(user, pass, host, db string) string {
	return fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", user, pass, host, db) // OK: connection string
}
