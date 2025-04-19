package golicenser

import (
	"github.com/go-viper/mapstructure/v2"
	"github.com/joshuasing/golicenser"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/v2/pkg/config"
	"github.com/golangci/golangci-lint/v2/pkg/goanalysis"
	"github.com/golangci/golangci-lint/v2/pkg/golinters/internal"
)

const (
	linterName = "golicenser"
	linterDesc = "Powerful license header linter"
)

func New(settings *config.GoLicenserSettings) *goanalysis.Linter {
	var conf golicenser.Config
	if settings != nil {
		var err error
		var yearMode golicenser.YearMode
		if ym := settings.Header.YearMode; ym != "" {
			yearMode, err = golicenser.ParseYearMode(ym)
			if err != nil {
				internal.LinterLogger.Fatalf("%s: parse year mode: %v", linterName, err)
			}
		}

		var commentStyle golicenser.CommentStyle
		if cs := settings.Header.CommentStyle; cs != "" {
			commentStyle, err = golicenser.ParseCommentStyle(cs)
			if err != nil {
				internal.LinterLogger.Fatalf("%s: parse comment style: %v", linterName, err)
			}
		}

		vars := make(map[string]*golicenser.Var, len(settings.Header.Variables))
		for k, v := range settings.Header.Variables {
			if s, ok := v.(string); ok {
				vars[k] = &golicenser.Var{Value: s}
				continue
			}

			var glVar config.GoLicenserVar
			if err := mapstructure.Decode(v, &glVar); err != nil {
				internal.LinterLogger.Fatalf("%s: decode variable %s: %v", linterName, k, err)
			}
			vars[k] = &golicenser.Var{
				Value:  glVar.Value,
				Regexp: glVar.Regexp,
			}
		}

		conf = golicenser.Config{
			Header: golicenser.HeaderOpts{
				Template:      settings.Header.Template,
				Matcher:       settings.Header.Matcher,
				MatcherEscape: settings.Header.MatcherEscape,
				Author:        settings.Header.Author,
				AuthorRegexp:  settings.Header.AuthorRegexp,
				Variables:     vars,
				YearMode:      yearMode,
				CommentStyle:  commentStyle,
			},
			Exclude:                settings.Exclude,
			MaxConcurrent:          settings.MaxConcurrent,
			CopyrightHeaderMatcher: settings.CopyrightHeaderMatcher,
		}
	}

	if conf.Header.Template == "" || conf.Header.Author == "" {
		// User did not set template or author, disable golicenser.
		return goanalysis.NewLinter(linterName, linterDesc, nil, nil)
	}

	analyzer, err := golicenser.NewAnalyzer(conf)
	if err != nil {
		internal.LinterLogger.Fatalf("%s: create analyzer: %v", linterName, err)
	}

	return goanalysis.NewLinter(
		linterName,
		linterDesc,
		[]*analysis.Analyzer{analyzer},
		nil,
	).WithLoadMode(goanalysis.LoadModeSyntax)
}
