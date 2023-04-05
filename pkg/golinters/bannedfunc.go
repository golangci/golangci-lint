package golinters

import (
	"go/ast"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"
)

type bannedFunc struct {
	ban *config.BannedFuncSettings
}

// NewBannedFunc returns a new banned function linter.
// example:
//   linters-settings:
//     bannedfunc:
//       - (ioutil).WriteFile: "As of Go 1.16, this function simply calls os.WriteFile."
//       - (ioutil).ReadFile: "As of Go 1.16, this function simply calls os.ReadFile."
//       - (github.com/example/banned).New: "This function is deprecated"
func NewBannedFunc(ban *config.BannedFuncSettings) *goanalysis.Linter {
	bf := &bannedFunc{ban: ban}
	return goanalysis.NewLinter(
		"bannedfunc",
		"Checks for use of banned functions",
		[]*analysis.Analyzer{
			{
				Name:     "bannedfunc",
				Doc:      "Checks for use of banned functions",
				Requires: []*analysis.Analyzer{inspect.Analyzer},
				Run:      bf.Run,
			},
		},
		nil,
	).WithLoadMode(goanalysis.LoadModeSyntax)
}

// Run runs the banned function linter.
func (bf *bannedFunc) Run(pass *analysis.Pass) (interface{}, error) {
	var (
		confMap = bf.parseBannedFunc()
		usedMap = make(map[string]map[string]string)
	)
	for _, imp := range pass.Pkg.Imports() {
		if conf, ok := confMap[imp.Path()]; ok {
			usedMap[imp.Path()] = conf
		}
	}

	for _, file := range pass.Files {
		ast.Inspect(file, func(n ast.Node) bool {
			selector, ok := n.(*ast.SelectorExpr)
			if !ok {
				return true
			}
			ident, ok := selector.X.(*ast.Ident)
			if !ok {
				return true
			}
			m := usedMap[ident.Name]
			if m == nil {
				return true
			}
			value, ok := m[selector.Sel.Name]
			if !ok {
				return true
			}
			pass.Reportf(n.Pos(), value)
			return true
		})
	}
	return nil, nil
}

// parseBannedFunc parses the banned function configuration.
// return: map[import]map[func]tips
// example:
// 	{
// 		"ioutil": {
// 			"WriteFile": "As of Go 1.16, this function simply calls os.WriteFile.",
// 			"ReadFile":"As of Go 1.16, this function simply calls os.ReadFile.",
// 		},
// 		"github.com/example/banned": {
// 			"New": "This function is deprecated",
// 		},
// 	}
func (bf *bannedFunc) parseBannedFunc() map[string]map[string]string {
	confMap := make(map[string]map[string]string, len(bf.ban.Funcs))
	for f, tips := range bf.ban.Funcs {
		first, last := strings.Index(f, "("), strings.Index(f, ")")
		if first < 0 || last <= 0 || first > last || first+1 == last {
			continue
		}
		var (
			importName = f[first+1 : last]
			funcName   = f[last+2:]
		)
		if importName == "" || funcName == "" {
			continue
		}
		if conf, ok := confMap[importName]; ok {
			conf[funcName] = tips
		} else {
			confMap[importName] = map[string]string{funcName: tips}
		}
	}
	return confMap
}
