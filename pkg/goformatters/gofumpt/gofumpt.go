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

func New(cfg config.GofumptSettings, goVersion string) *Formatter {
	return &Formatter{
		options: gofumpt.Options{
			LangVersion: getLangVersion(goVersion),
			ModulePath:  cfg.ModulePath,
			ExtraRules:  cfg.ExtraRules,
		},
	}
}

func (*Formatter) Name() string {
	return Name
}

func (f *Formatter) Format(_ string, src []byte) ([]byte, error) {
	return gofumpt.Source(src, f.options)
}

// modified copy of pkg/golinters/gofumpt/gofumpt.go
func getLangVersion(v string) string {
	if v == "" {
		// TODO: defaults to "1.15", in the future (v2) must be removed.
		return "go1.15"
	}

	return "go" + strings.TrimPrefix(v, "go")
}
