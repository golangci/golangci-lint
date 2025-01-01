package misspell

import (
	"fmt"
	"go/ast"
	"go/token"
	"strings"
	"unicode"

	"github.com/golangci/misspell"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/goanalysis"
	"github.com/golangci/golangci-lint/pkg/lint/linter"
)

const linterName = "misspell"

func New(settings *config.MisspellSettings) *goanalysis.Linter {
	analyzer := &analysis.Analyzer{
		Name: linterName,
		Doc:  goanalysis.TheOnlyanalyzerDoc,
		Run:  goanalysis.DummyRun,
	}

	return goanalysis.NewLinter(
		linterName,
		"Finds commonly misspelled English words",
		[]*analysis.Analyzer{analyzer},
		nil,
	).WithContextSetter(func(lintCtx *linter.Context) {
		replacer, ruleErr := createMisspellReplacer(settings)

		analyzer.Run = func(pass *analysis.Pass) (any, error) {
			if ruleErr != nil {
				return nil, ruleErr
			}

			err := runMisspell(lintCtx, pass, replacer, settings.Mode)
			if err != nil {
				return nil, err
			}

			return nil, nil
		}
	}).WithLoadMode(goanalysis.LoadModeSyntax)
}

func runMisspell(lintCtx *linter.Context, pass *analysis.Pass, replacer *misspell.Replacer, mode string) error {
	for _, file := range pass.Files {
		err := runMisspellOnFile(lintCtx, pass, file, replacer, mode)
		if err != nil {
			return err
		}
	}

	return nil
}

func createMisspellReplacer(settings *config.MisspellSettings) (*misspell.Replacer, error) {
	replacer := &misspell.Replacer{
		Replacements: misspell.DictMain,
	}

	// Figure out regional variations
	switch strings.ToUpper(settings.Locale) {
	case "":
		// nothing
	case "US":
		replacer.AddRuleList(misspell.DictAmerican)
	case "UK", "GB":
		replacer.AddRuleList(misspell.DictBritish)
	case "NZ", "AU", "CA":
		return nil, fmt.Errorf("unknown locale: %q", settings.Locale)
	}

	err := appendExtraWords(replacer, settings.ExtraWords)
	if err != nil {
		return nil, fmt.Errorf("process extra words: %w", err)
	}

	if len(settings.IgnoreWords) != 0 {
		replacer.RemoveRule(settings.IgnoreWords)
	}

	// It can panic.
	replacer.Compile()

	return replacer, nil
}

func runMisspellOnFile(lintCtx *linter.Context, pass *analysis.Pass, file *ast.File, replacer *misspell.Replacer, mode string) error {
	position, isGoFile := goanalysis.GetGoFilePosition(pass, file)
	if !isGoFile {
		return nil
	}

	fileContent, err := lintCtx.FileCache.GetFileBytes(position.Filename)
	if err != nil {
		return fmt.Errorf("can't get file %s contents: %w", position.Filename, err)
	}

	// `r.ReplaceGo` doesn't find issues inside strings: it searches only inside comments.
	// `r.Replace` searches all words: it treats input as a plain text.
	// The standalone misspell tool uses `r.Replace` by default.
	var replace func(input string) (string, []misspell.Diff)
	switch strings.ToLower(mode) {
	case "restricted":
		replace = replacer.ReplaceGo
	default:
		replace = replacer.Replace
	}

	f := pass.Fset.File(file.Pos())

	_, diffs := replace(string(fileContent))

	for _, diff := range diffs {
		text := fmt.Sprintf("`%s` is a misspelling of `%s`", diff.Original, diff.Corrected)

		start := f.LineStart(diff.Line) + token.Pos(diff.Column)
		end := f.LineStart(diff.Line) + token.Pos(diff.Column+len(diff.Original))

		pass.Report(analysis.Diagnostic{
			Pos:     start,
			End:     end,
			Message: text,
			SuggestedFixes: []analysis.SuggestedFix{{
				TextEdits: []analysis.TextEdit{{
					Pos:     start,
					End:     end,
					NewText: []byte(diff.Corrected),
				}},
			}},
		})
	}

	return nil
}

func appendExtraWords(replacer *misspell.Replacer, extraWords []config.MisspellExtraWords) error {
	if len(extraWords) == 0 {
		return nil
	}

	extra := make([]string, 0, len(extraWords)*2)

	for _, word := range extraWords {
		if word.Typo == "" || word.Correction == "" {
			return fmt.Errorf("typo (%q) and correction (%q) fields should not be empty", word.Typo, word.Correction)
		}

		if strings.ContainsFunc(word.Typo, func(r rune) bool { return !unicode.IsLetter(r) }) {
			return fmt.Errorf("the word %q in the 'typo' field should only contain letters", word.Typo)
		}
		if strings.ContainsFunc(word.Correction, func(r rune) bool { return !unicode.IsLetter(r) }) {
			return fmt.Errorf("the word %q in the 'correction' field should only contain letters", word.Correction)
		}

		extra = append(extra, strings.ToLower(word.Typo), strings.ToLower(word.Correction))
	}

	replacer.AddRuleList(extra)

	return nil
}
