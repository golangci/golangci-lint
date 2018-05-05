package golinters

import "github.com/golangci/golangci-lint/pkg"

const pathLineColMessage = `^(?P<path>.*?\.go):(?P<line>\d+):(?P<col>\d+):\s*(?P<message>.*)$`
const pathLineMessage = `^(?P<path>.*?\.go):(?P<line>\d+):\s*(?P<message>.*)$`

var errCheck = newLinter("errcheck",
	newLinterConfig(
		"Error return value is not checked",
		pathLineColMessage,
		"\\.Close()", // It's annoying and not critical error to ignore Close() errors),
	),
)

var golint = newLinter("golint", newLinterConfig("", pathLineColMessage, ""))
var govet = newLinter("govet", newLinterConfig("", pathLineMessage, "", "--no-recurse"))

func GetSupportedLinters() []linters.Linter {
	return []linters.Linter{errCheck}
}
