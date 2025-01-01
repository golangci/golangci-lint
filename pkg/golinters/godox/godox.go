package godox

import (
	"go/token"
	"strings"

	"github.com/matoous/godox"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/goanalysis"
)

const linterName = "godox"

func New(settings *config.GodoxSettings) *goanalysis.Linter {
	analyzer := &analysis.Analyzer{
		Name: linterName,
		Doc:  goanalysis.TheOnlyanalyzerDoc,
		Run: func(pass *analysis.Pass) (any, error) {
			runGodox(pass, settings)

			return nil, nil
		},
	}

	return goanalysis.NewLinter(
		linterName,
		"Tool for detection of FIXME, TODO and other comment keywords",
		[]*analysis.Analyzer{analyzer},
		nil,
	).WithLoadMode(goanalysis.LoadModeSyntax)
}

func runGodox(pass *analysis.Pass, settings *config.GodoxSettings) {
	for _, file := range pass.Files {
		position, isGoFile := goanalysis.GetGoFilePosition(pass, file)
		if !isGoFile {
			continue
		}

		messages := godox.Run(file, pass.Fset, settings.Keywords...)
		if len(messages) == 0 {
			continue
		}

		nonAdjPosition := pass.Fset.PositionFor(file.Pos(), false)

		ft := pass.Fset.File(file.Pos())

		for _, i := range messages {
			pass.Report(analysis.Diagnostic{
				Pos:     ft.LineStart(goanalysis.AdjustPos(i.Pos.Line, nonAdjPosition.Line, position.Line)) + token.Pos(i.Pos.Column),
				Message: strings.TrimRight(i.Message, "\n"),
			})
		}
	}
}
