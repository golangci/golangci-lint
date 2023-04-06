package golinters

import (
	"github.com/demoManito/bannedfunc"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"
)

// NewBannedFunc returns a new bannedfunc linter.
func NewBannedFunc(ban *config.BannedFuncSettings) *goanalysis.Linter {
	return goanalysis.NewLinter(
		bannedfunc.Name,
		bannedfunc.Doc,
		[]*analysis.Analyzer{
			{
				Name:     "bannedfunc",
				Doc:      "Checks for use of banned functions",
				Requires: []*analysis.Analyzer{inspect.Analyzer},
				Run: func(pass *analysis.Pass) (interface{}, error) {
					msgs := bannedfunc.NewLinter(ban.Funcs, pass.Pkg, pass.Files).Run()
					for _, msg := range msgs {
						pass.Reportf(msg.Pos, msg.Tips)
					}
					return nil, nil
				},
			},
		},
		nil,
	).WithLoadMode(goanalysis.LoadModeSyntax)
}
