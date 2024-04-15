package sloglint

import (
	"go-simpler.org/sloglint"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/goanalysis"
)

func New(settings *config.SlogLintSettings) *goanalysis.Linter {
	var opts *sloglint.Options
	if settings != nil {
		opts = &sloglint.Options{
			NoMixedArgs:    settings.NoMixedArgs,
			KVOnly:         settings.KVOnly,
			AttrOnly:       settings.AttrOnly,
			NoGlobal:       settings.NoGlobal,
			ContextOnly:    settings.Context,
			StaticMsg:      settings.StaticMsg,
			NoRawKeys:      settings.NoRawKeys,
			KeyNamingCase:  settings.KeyNamingCase,
			ArgsOnSepLines: settings.ArgsOnSepLines,
		}
	}

	a := sloglint.New(opts)

	return goanalysis.
		NewLinter(a.Name, a.Doc, []*analysis.Analyzer{a}, nil).
		WithLoadMode(goanalysis.LoadModeTypesInfo)
}
