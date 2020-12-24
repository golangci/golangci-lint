package golinters

import (
	"bufio"
	"bytes"
	"go/token"
	"io/ioutil"
	"strings"
	"sync"

	"github.com/ashanbrown/gofmts/pkg/gofmts"
	"github.com/pkg/errors"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"
	"github.com/golangci/golangci-lint/pkg/lint/linter"
	"github.com/golangci/golangci-lint/pkg/result"
)

const gofmtsName = "gofmts"

func NewGofmts() *goanalysis.Linter {
	var mu sync.Mutex
	var resIssues []goanalysis.Issue

	analyzer := &analysis.Analyzer{
		Name: gofmtsName,
		Doc:  goanalysis.TheOnlyanalyzerDoc,
	}
	return goanalysis.NewLinter(
		gofmtsName,
		"Gofmts checks formatting of marked strings and sorting of expression groups.",
		[]*analysis.Analyzer{analyzer},
		nil,
	).WithContextSetter(func(lintCtx *linter.Context) {
		analyzer.Run = func(pass *analysis.Pass) (interface{}, error) {
			sorter := gofmts.NewSorter()
			formatter := gofmts.NewFormatter()
			for _, f := range pass.Files {
				pos := pass.Fset.PositionFor(f.Pos(), false)
				src, err := ioutil.ReadFile(pos.Filename)
				if err != nil {
					return nil, errors.Wrapf(err, "%s linter failed to read file %q", gofmtsName, f.Name.String())
				}

				sortIssues, err := sorter.Run(pass.Fset, f)
				if err != nil {
					return nil, errors.Wrapf(err, "%s linter failed on file %q", gofmtsName, f.Name.String())
				}
				mu.Lock()
				resIssues = append(resIssues, processGofmtsIssues(pass, src, sortIssues)...)
				mu.Unlock()
				
				formatIssues, err := formatter.Run(src, pass.Fset, f)
				if err != nil {
					return nil, errors.Wrapf(err, "%s linter failed on file %q", gofmtsName, f.Name.String())
				}

				mu.Lock()
				resIssues = append(resIssues, processGofmtsIssues(pass, src, formatIssues)...)
				mu.Unlock()
			}
			
			return nil, nil
		}
	}).WithIssuesReporter(func(*linter.Context) []goanalysis.Issue {
		return resIssues
	}).WithLoadMode(goanalysis.LoadModeSyntax)
}

func processGofmtsIssues(pass *analysis.Pass, src []byte, issues []gofmts.Issue) []goanalysis.Issue {
	var results []goanalysis.Issue
	for _, i := range issues {
		issue := goanalysis.NewIssue(&result.Issue{
			Pos:        i.Position(),
			Text:       i.Details(),
			FromLinter: gofmtsName,
		}, pass)

		if i, hasReplacement := i.(gofmts.IssueWithReplacement); hasReplacement {
			// find the end position in the file 
			startPosition := i.Position()
			var end token.Pos
			pass.Fset.Iterate(func(f *token.File) bool {
				if f.Name() == startPosition.Filename {
					end = f.Pos(startPosition.Offset + i.Length())
					return false
				}
				return true
			})
			endPosition := pass.Fset.Position(end)

			if i.Position().Line == endPosition.Line {
				issue.Replacement = &result.Replacement{
					Inline: &result.InlineFix{
						StartCol:  i.Position().Column - 1,
						Length:    i.Length(),
						NewString: i.Replacement(),
					},
				}
			} else {
				issue.LineRange = &result.Range{
					From: startPosition.Line,
					To:   endPosition.Line,
				}

				var replacementLines []string
				scanner := bufio.NewScanner(strings.NewReader(i.Replacement()))
				for scanner.Scan() {
					replacementLines = append(replacementLines, scanner.Text())
				}

				// add rest of first and last lines to non-sort issues
				// sort issues return a replacement for the entire line but others don't
				if _, isSort := i.(gofmts.SortIssue); !isSort {
					issue.LineRange.To = endPosition.Line + 1
					
					// expand the replacement to update entire lines
					firstLineRest := string(src[startPosition.Offset-(startPosition.Column-1) : startPosition.Offset])

					// find the rest of the last line
					endScanner := bufio.NewScanner(bytes.NewReader(src[endPosition.Offset:]))
					_ = endScanner.Scan()
					lastLineRest := i.Replacement() + endScanner.Text()

					replacementLines[0] = firstLineRest + replacementLines[0]
					replacementLines[len(replacementLines)-1] += lastLineRest
				}

				issue.Replacement = &result.Replacement{
					NewLines: replacementLines,
				}
			}
		}
		results = append(results, issue)
	}
	return results
}