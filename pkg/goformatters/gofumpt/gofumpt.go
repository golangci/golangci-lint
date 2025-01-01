package gofumpt

import (
	"strings"

	gofumpt "mvdan.cc/gofumpt/format"

	"github.com/golangci/golangci-lint/pkg/config"
)

const Name = "gofumpt"

type Formatter struct {
	options gofumpt.Options
}

func New(settings *config.GofumptSettings, goVersion string) *Formatter {
	var options gofumpt.Options

	if settings != nil {
		options = gofumpt.Options{
			LangVersion: getLangVersion(goVersion),
			ModulePath:  settings.ModulePath,
			ExtraRules:  settings.ExtraRules,
		}
	}

	return &Formatter{options: options}
}

func (*Formatter) Name() string {
	return Name
}

func (f *Formatter) Format(_ string, src []byte) ([]byte, error) {
	return gofumpt.Source(src, f.options)
}

func getLangVersion(v string) string {
	if v == "" {
		// TODO: defaults to "1.15", in the future (v2) must be removed.
		return "go1.15"
	}

	return "go" + strings.TrimPrefix(v, "go")
}
