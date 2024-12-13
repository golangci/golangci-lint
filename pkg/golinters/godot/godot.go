package godot

import (
	"cmp"
	"go/token"

	"github.com/tetafro/godot"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/goanalysis"
)

const linterName = "godot"

func New(settings *config.GodotSettings) *goanalysis.Linter {
	var dotSettings godot.Settings

	if settings != nil {
		dotSettings = godot.Settings{
			Scope:   godot.Scope(settings.Scope),
			Exclude: settings.Exclude,
			Period:  settings.Period,
			Capital: settings.Capital,
		}

		// Convert deprecated setting
		if settings.CheckAll {
			dotSettings.Scope = godot.AllScope
		}
	}

	dotSettings.Scope = cmp.Or(dotSettings.Scope, godot.DeclScope)

	analyzer := &analysis.Analyzer{
		Name: linterName,
		Doc:  goanalysis.TheOnlyanalyzerDoc,
		Run: func(pass *analysis.Pass) (any, error) {
			err := runGodot(pass, dotSettings)
			if err != nil {
				return nil, err
			}

			return nil, nil
		},
	}

	return goanalysis.NewLinter(
		linterName,
		"Check if comments end in a period",
		[]*analysis.Analyzer{analyzer},
		nil,
	).WithLoadMode(goanalysis.LoadModeSyntax)
}

func runGodot(pass *analysis.Pass, settings godot.Settings) error {
	for _, file := range pass.Files {
		iss, err := godot.Run(file, pass.Fset, settings)
		if err != nil {
			return err
		}

		for _, i := range iss {
			f := pass.Fset.File(file.Pos())

			pos := f.Pos(i.Pos.Offset)

			var end token.Pos
			if i.Pos.Line == f.LineCount() {
				// missing newline at the end of the file
				end = f.Pos(f.Size())
			} else {
				end = f.LineStart(i.Pos.Line+1) - token.Pos(1)
			}

			pass.Report(analysis.Diagnostic{
				Pos:     pos,
				End:     end,
				Message: i.Message,
				SuggestedFixes: []analysis.SuggestedFix{{
					TextEdits: []analysis.TextEdit{{
						Pos:     pos,
						End:     end,
						NewText: []byte(i.Replacement),
					}},
				}},
			})
		}
	}

	return nil
}
