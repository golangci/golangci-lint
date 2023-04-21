package golinters

import (
	"go.tmz.dev/musttag"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"
)

func NewMustTag(setting *config.MustTagSettings) *goanalysis.Linter {
	var funcs []musttag.Func

	if setting != nil {
		for _, fn := range setting.Functions {
			funcs = append(funcs, musttag.Func{
				Name:   fn.Name,
				Tag:    fn.Tag,
				ArgPos: fn.ArgPos,
			})
		}
	}

	a := musttag.New(funcs...)

	return goanalysis.
		NewLinter(a.Name, a.Doc, []*analysis.Analyzer{a}, nil).
		WithLoadMode(goanalysis.LoadModeTypesInfo)
}
