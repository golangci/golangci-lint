package golinters

import (
	"go/token"

	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"
)

const (
	multiImportName = `multiimport`
	multiImportDesc = `Finds files where packages are imported more than once under different aliases.`
)

func NewMultiImport() *goanalysis.Linter {
	a := &analysis.Analyzer{
		Name: multiImportName,
		Doc:  multiImportDesc,
		Run:  runMultiImport,
	}
	return goanalysis.
		NewLinter(multiImportName, multiImportDesc, []*analysis.Analyzer{a}, nil).
		WithLoadMode(goanalysis.LoadModeSyntax)
}

func runMultiImport(pass *analysis.Pass) (interface{}, error) {
	for _, f := range pass.Files {
		importToPathPos := make(map[string][]token.Pos, len(f.Imports))

		for _, decl := range f.Imports {
			importToPathPos[decl.Path.Value] = append(importToPathPos[decl.Path.Value], decl.Path.ValuePos)
		}

		for _, positions := range importToPathPos {
			if len(positions) <= 1 {
				continue
			}
			for _, pos := range positions {
				pass.Reportf(pos, `import appears multiple times under different aliases (%s)`, multiImportName)
			}
		}
	}
	return nil, nil
}
