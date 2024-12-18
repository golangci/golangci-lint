package goheader

import (
	"go/token"
	"strings"

	goheader "github.com/denis-tingaikin/go-header"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/goanalysis"
)

const linterName = "goheader"

func New(settings *config.GoHeaderSettings) *goanalysis.Linter {
	conf := &goheader.Configuration{}
	if settings != nil {
		conf = &goheader.Configuration{
			Values:       settings.Values,
			Template:     settings.Template,
			TemplatePath: settings.TemplatePath,
		}
	}

	analyzer := &analysis.Analyzer{
		Name: linterName,
		Doc:  goanalysis.TheOnlyanalyzerDoc,
		Run: func(pass *analysis.Pass) (any, error) {
			err := runGoHeader(pass, conf)
			if err != nil {
				return nil, err
			}

			return nil, nil
		},
	}

	return goanalysis.NewLinter(
		linterName,
		"Checks if file header matches to pattern",
		[]*analysis.Analyzer{analyzer},
		nil,
	).WithLoadMode(goanalysis.LoadModeSyntax)
}

func runGoHeader(pass *analysis.Pass, conf *goheader.Configuration) error {
	if conf.TemplatePath == "" && conf.Template == "" {
		// User did not pass template, so then do not run go-header linter
		return nil
	}

	template, err := conf.GetTemplate()
	if err != nil {
		return err
	}

	values, err := conf.GetValues()
	if err != nil {
		return err
	}

	a := goheader.New(goheader.WithTemplate(template), goheader.WithValues(values))

	for _, file := range pass.Files {
		position := goanalysis.GetFilePosition(pass, file)

		issue := a.Analyze(&goheader.Target{File: file, Path: position.Filename})
		if issue == nil {
			continue
		}

		f := pass.Fset.File(file.Pos())

		commentLine := 1

		// Inspired by https://github.com/denis-tingaikin/go-header/blob/4c75a6a2332f025705325d6c71fff4616aedf48f/analyzer.go#L85-L92
		if len(file.Comments) > 0 && file.Comments[0].Pos() < file.Package {
			commentLine = goanalysis.GetFilePositionFor(pass.Fset, file.Comments[0].Pos()).Line
		}

		start := f.LineStart(commentLine)

		diag := analysis.Diagnostic{
			Pos:     start,
			Message: issue.Message(),
		}

		if fix := issue.Fix(); fix != nil {
			end := len(fix.Actual)
			for _, s := range fix.Actual {
				end += len(s)
			}

			diag.SuggestedFixes = []analysis.SuggestedFix{{
				TextEdits: []analysis.TextEdit{{
					Pos:     start,
					End:     start + token.Pos(end),
					NewText: []byte(strings.Join(fix.Expected, "\n") + "\n"),
				}},
			}}
		}

		pass.Report(diag)
	}

	return nil
}
