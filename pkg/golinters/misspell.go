package golinters

import (
	"context"
	"fmt"
	"go/token"
	"strings"

	"github.com/golangci/misspell"

	"github.com/golangci/golangci-lint/pkg/lint/linter"
	"github.com/golangci/golangci-lint/pkg/result"
)

type Misspell struct{}

func NewMisspell() *Misspell {
	return &Misspell{}
}

func (Misspell) Name() string {
	return "misspell"
}

func (Misspell) Desc() string {
	return "Finds commonly misspelled English words in comments"
}

func (lint Misspell) Run(ctx context.Context, lintCtx *linter.Context) ([]result.Issue, error) {
	r := misspell.Replacer{
		Replacements: misspell.DictMain,
	}

	// Figure out regional variations
	settings := lintCtx.Settings().Misspell
	locale := settings.Locale
	switch strings.ToUpper(locale) {
	case "":
		// nothing
	case "US":
		r.AddRuleList(misspell.DictAmerican)
	case "UK", "GB":
		r.AddRuleList(misspell.DictBritish)
	case "NZ", "AU", "CA":
		return nil, fmt.Errorf("unknown locale: %q", locale)
	}

	if len(settings.IgnoreWords) != 0 {
		r.RemoveRule(settings.IgnoreWords)
	}

	r.Compile()

	var res []result.Issue
	for _, f := range getAllFileNames(lintCtx) {
		issues, err := lint.runOnFile(f, &r, lintCtx)
		if err != nil {
			return nil, err
		}
		res = append(res, issues...)
	}

	return res, nil
}

func (lint Misspell) runOnFile(fileName string, r *misspell.Replacer, lintCtx *linter.Context) ([]result.Issue, error) {
	var res []result.Issue
	fileContent, err := lintCtx.FileCache.GetFileBytes(fileName)
	if err != nil {
		return nil, fmt.Errorf("can't get file %s contents: %s", fileName, err)
	}

	// use r.Replace, not r.ReplaceGo because r.ReplaceGo doesn't find
	// issues inside strings: it searches only inside comments. r.Replace
	// searches all words: it treats input as a plain text. A standalone misspell
	// tool uses r.Replace by default.
	_, diffs := r.Replace(string(fileContent))
	for _, diff := range diffs {
		text := fmt.Sprintf("`%s` is a misspelling of `%s`", diff.Original, diff.Corrected)
		pos := token.Position{
			Filename: fileName,
			Line:     diff.Line,
			Column:   diff.Column + 1,
		}
		replacement := &result.Replacement{
			Inline: &result.InlineFix{
				StartCol:  diff.Column,
				Length:    len(diff.Original),
				NewString: diff.Corrected,
			},
		}

		res = append(res, result.Issue{
			Pos:         pos,
			Text:        text,
			FromLinter:  lint.Name(),
			Replacement: replacement,
		})
	}

	return res, nil
}
