package golicenser

import (
	"fmt"
	"os"
	"strings"

	"github.com/joshuasing/golicenser"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/v2/pkg/config"
	"github.com/golangci/golangci-lint/v2/pkg/goanalysis"
	"github.com/golangci/golangci-lint/v2/pkg/golinters/internal"
)

const linterName = "golicenser"

func New(settings *config.GoLicenserSettings, replacer *strings.Replacer) *goanalysis.Linter {
	conf, err := createConfig(settings, replacer)
	if err != nil {
		internal.LinterLogger.Fatalf("%s: parse year mode: %v", linterName, err)
	}

	analyzer, err := golicenser.NewAnalyzer(conf)
	if err != nil {
		internal.LinterLogger.Fatalf("%s: create analyzer: %v", linterName, err)
	}

	return goanalysis.NewLinter(
		linterName,
		"Powerful license header linter",
		[]*analysis.Analyzer{analyzer},
		nil,
	).WithLoadMode(goanalysis.LoadModeSyntax)
}

func createConfig(settings *config.GoLicenserSettings, replacer *strings.Replacer) (golicenser.Config, error) {
	if settings == nil {
		return golicenser.Config{}, nil
	}

	header := golicenser.HeaderOpts{
		Matcher:       settings.Header.Matcher,
		MatcherEscape: settings.Header.MatcherEscape,
		Author:        settings.Header.Author,
		AuthorRegexp:  settings.Header.AuthorRegexp,
		Template:      settings.Header.Template,
		Variables:     make(map[string]*golicenser.Var, len(settings.Header.Variables)),
	}

	// Template from config takes priority over template from 'template-path'.
	if header.Template == "" && settings.Header.TemplatePath != "" {
		b, err := os.ReadFile(replacer.Replace(settings.Header.TemplatePath))
		if err != nil {
			return golicenser.Config{}, fmt.Errorf("read the template file: %w", err)
		}

		// Use template from a file (trim newline from the end of file)
		header.Template = strings.TrimSuffix(strings.TrimSuffix(string(b), "\n"), "\r")
	}

	var err error
	if settings.Header.YearMode != "" {
		header.YearMode, err = golicenser.ParseYearMode(settings.Header.YearMode)
		if err != nil {
			return golicenser.Config{}, fmt.Errorf("parse year mode: %w", err)
		}
	}

	if settings.Header.CommentStyle != "" {
		header.CommentStyle, err = golicenser.ParseCommentStyle(settings.Header.CommentStyle)
		if err != nil {
			return golicenser.Config{}, fmt.Errorf("parse comment style: %w", err)
		}
	}

	for k, v := range settings.Header.Variables {
		header.Variables[k] = &golicenser.Var{
			Value:  v.Value,
			Regexp: v.Regexp,
		}
	}

	return golicenser.Config{
		Header:                 header,
		CopyrightHeaderMatcher: settings.CopyrightHeaderMatcher,

		// NOTE(ldez): not wanted for now.
		// And an empty slice is used because of a wrong default inside the linter.
		Exclude: []string{},
		// NOTE(ldez): golangci-lint already handles concurrency.
		MaxConcurrent: 1,
	}, nil
}
