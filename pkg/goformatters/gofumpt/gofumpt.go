package gofumpt

import (
	"strings"

	gofumpt "mvdan.cc/gofumpt/format"

	"github.com/golangci/golangci-lint/v2/pkg/config"
	"github.com/golangci/golangci-lint/v2/pkg/goformatters/internal"
)

const Name = "gofumpt"

type Formatter struct {
	options gofumpt.Options
}

func New(settings *config.GoFumptSettings, goVersion string) *Formatter {
	var options gofumpt.Options

	if settings != nil {
		if settings.ExtraRules {
			internal.FormatterLogger.Warnf("gofumpt: `extra-rules` is deprecated, please use `extra.group-params` instead.")
		}

		options = gofumpt.Options{
			LangVersion: getLangVersion(goVersion),
			ModulePath:  settings.ModulePath,
			ExtraRules:  settings.ExtraRules,
			Extra: gofumpt.Extra{
				GroupParams:   settings.Extra.GroupParams,
				ClotheReturns: settings.Extra.ClotheReturns,
			},
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
	return "go" + strings.TrimPrefix(v, "go")
}
