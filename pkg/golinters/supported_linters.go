package golinters

import "github.com/golangci/golangci-lint/pkg"

const pathLineColMessage = `^(?P<path>.*?\.go):(?P<line>\d+):(?P<col>\d+):\s*(?P<message>.*)$`
const pathLineMessage = `^(?P<path>.*?\.go):(?P<line>\d+):\s*(?P<message>.*)$`

func GetSupportedLinters() []pkg.Linter {
	return []pkg.Linter{govet{}, errcheck{}, golint{}}
}
