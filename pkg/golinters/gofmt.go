package golinters

import (
	"bytes"
	"context"
	"fmt"
	"go/token"
	"strings"

	gofmtAPI "github.com/golangci/gofmt/gofmt"
	goimportsAPI "github.com/golangci/gofmt/goimports"
	"github.com/pkg/errors"
	"golang.org/x/tools/imports"
	diffpkg "sourcegraph.com/sourcegraph/go-diff/diff"

	"github.com/golangci/golangci-lint/pkg/lint/linter"
	"github.com/golangci/golangci-lint/pkg/logutils"
	"github.com/golangci/golangci-lint/pkg/result"
)

type Gofmt struct {
	UseGoimports bool
}

func (g Gofmt) Name() string {
	if g.UseGoimports {
		return "goimports"
	}

	return "gofmt"
}

func (g Gofmt) Desc() string {
	if g.UseGoimports {
		return "Goimports does everything that gofmt does. Additionally it checks unused imports"
	}

	return "Gofmt checks whether code was gofmt-ed. By default " +
		"this tool runs with -s option to check for code simplification"
}

type Change struct {
	LineRange   result.Range
	Replacement result.Replacement
}

type diffLineType string

const (
	diffLineAdded    diffLineType = "added"
	diffLineOriginal diffLineType = "original"
	diffLineDeleted  diffLineType = "deleted"
)

type diffLine struct {
	originalNumber int // 1-based original line number
	typ            diffLineType
	data           string // "+" or "-" stripped line
}

type hunkChangesParser struct {
	// needed because we merge currently added lines with the last original line
	lastOriginalLine *diffLine

	// if the first line of diff is an adding we save all additions to replacementLinesToPrepend
	replacementLinesToPrepend []string

	log logutils.Log

	lines []diffLine

	ret []Change
}

func (p *hunkChangesParser) parseDiffLines(h *diffpkg.Hunk) {
	lines := bytes.Split(h.Body, []byte{'\n'})
	currentOriginalLineNumer := int(h.OrigStartLine)
	var ret []diffLine

	for i, line := range lines {
		dl := diffLine{
			originalNumber: currentOriginalLineNumer,
		}

		lineStr := string(line)

		//nolint:gocritic
		if strings.HasPrefix(lineStr, "-") {
			dl.typ = diffLineDeleted
			dl.data = strings.TrimPrefix(lineStr, "-")
			currentOriginalLineNumer++
		} else if strings.HasPrefix(lineStr, "+") {
			dl.typ = diffLineAdded
			dl.data = strings.TrimPrefix(lineStr, "+")
		} else {
			if i == len(lines)-1 && lineStr == "" {
				// handle last \n: don't add an empty original line
				break
			}

			dl.typ = diffLineOriginal
			dl.data = strings.TrimPrefix(lineStr, " ")
			currentOriginalLineNumer++
		}

		ret = append(ret, dl)
	}

	p.lines = ret
}

func (p *hunkChangesParser) handleOriginalLine(line diffLine, i *int) {
	if len(p.replacementLinesToPrepend) == 0 {
		p.lastOriginalLine = &line
		*i++
		return
	}

	// check following added lines for the case:
	// + added line 1
	// original line
	// + added line 2

	*i++
	var followingAddedLines []string
	for ; *i < len(p.lines) && p.lines[*i].typ == diffLineAdded; *i++ {
		followingAddedLines = append(followingAddedLines, p.lines[*i].data)
	}

	p.ret = append(p.ret, Change{
		LineRange: result.Range{
			From: line.originalNumber,
			To:   line.originalNumber,
		},
		Replacement: result.Replacement{
			NewLines: append(p.replacementLinesToPrepend, append([]string{line.data}, followingAddedLines...)...),
		},
	})
	p.replacementLinesToPrepend = nil
	p.lastOriginalLine = &line
}

func (p *hunkChangesParser) handleDeletedLines(deletedLines []diffLine, addedLines []string) {
	change := Change{
		LineRange: result.Range{
			From: deletedLines[0].originalNumber,
			To:   deletedLines[len(deletedLines)-1].originalNumber,
		},
	}

	if len(addedLines) != 0 {
		//nolint:gocritic
		change.Replacement.NewLines = append(p.replacementLinesToPrepend, addedLines...)
		if len(p.replacementLinesToPrepend) != 0 {
			p.replacementLinesToPrepend = nil
		}

		p.ret = append(p.ret, change)
		return
	}

	// delete-only change with possible prepending
	if len(p.replacementLinesToPrepend) != 0 {
		change.Replacement.NewLines = p.replacementLinesToPrepend
		p.replacementLinesToPrepend = nil
	} else {
		change.Replacement.NeedOnlyDelete = true
	}

	p.ret = append(p.ret, change)
}

func (p *hunkChangesParser) handleAddedOnlyLines(addedLines []string) {
	if p.lastOriginalLine == nil {
		// the first line is added; the diff looks like:
		// 1. + ...
		// 2. - ...
		// or
		// 1. + ...
		// 2. ...

		p.replacementLinesToPrepend = addedLines
		return
	}

	// add-only change merged into the last original line with possible prepending
	p.ret = append(p.ret, Change{
		LineRange: result.Range{
			From: p.lastOriginalLine.originalNumber,
			To:   p.lastOriginalLine.originalNumber,
		},
		Replacement: result.Replacement{
			NewLines: append(p.replacementLinesToPrepend, append([]string{p.lastOriginalLine.data}, addedLines...)...),
		},
	})
	p.replacementLinesToPrepend = nil
}

func (p *hunkChangesParser) parse(h *diffpkg.Hunk) []Change {
	p.parseDiffLines(h)

	for i := 0; i < len(p.lines); {
		line := p.lines[i]
		if line.typ == diffLineOriginal {
			p.handleOriginalLine(line, &i) //nolint:scopelint
			continue
		}

		var deletedLines []diffLine
		for ; i < len(p.lines) && p.lines[i].typ == diffLineDeleted; i++ {
			deletedLines = append(deletedLines, p.lines[i])
		}

		var addedLines []string
		for ; i < len(p.lines) && p.lines[i].typ == diffLineAdded; i++ {
			addedLines = append(addedLines, p.lines[i].data)
		}

		if len(deletedLines) != 0 {
			p.handleDeletedLines(deletedLines, addedLines)
			continue
		}

		// no deletions, only addings
		p.handleAddedOnlyLines(addedLines)
	}

	if len(p.replacementLinesToPrepend) != 0 {
		p.log.Infof("The diff contains only additions: no original or deleted lines: %#v", p.lines)
		return nil
	}

	return p.ret
}

func (g Gofmt) extractIssuesFromPatch(patch string, log logutils.Log) ([]result.Issue, error) {
	diffs, err := diffpkg.ParseMultiFileDiff([]byte(patch))
	if err != nil {
		return nil, errors.Wrap(err, "can't parse patch")
	}

	if len(diffs) == 0 {
		return nil, fmt.Errorf("got no diffs from patch parser: %v", diffs)
	}

	issues := []result.Issue{}
	for _, d := range diffs {
		if len(d.Hunks) == 0 {
			log.Warnf("Got no hunks in diff %+v", d)
			continue
		}

		for _, hunk := range d.Hunks {
			text := "File is not `gofmt`-ed with `-s`"
			if g.UseGoimports {
				text = "File is not `goimports`-ed"
			}
			p := hunkChangesParser{
				log: log,
			}
			changes := p.parse(hunk)
			for _, change := range changes {
				change := change // fix scope
				i := result.Issue{
					FromLinter: g.Name(),
					Pos: token.Position{
						Filename: d.NewName,
						Line:     change.LineRange.From,
					},
					Text:        text,
					Replacement: &change.Replacement,
				}
				if change.LineRange.From != change.LineRange.To {
					i.LineRange = &change.LineRange
				}

				issues = append(issues, i)
			}
		}
	}

	return issues, nil
}

func (g Gofmt) Run(ctx context.Context, lintCtx *linter.Context) ([]result.Issue, error) {
	var issues []result.Issue

	for _, f := range getAllFileNames(lintCtx) {
		var diff []byte
		var err error
		if g.UseGoimports {
			imports.LocalPrefix = lintCtx.Settings().Goimports.LocalPrefixes
			diff, err = goimportsAPI.Run(f)
		} else {
			diff, err = gofmtAPI.Run(f, lintCtx.Settings().Gofmt.Simplify)
		}
		if err != nil { // TODO: skip
			return nil, err
		}
		if diff == nil {
			continue
		}

		is, err := g.extractIssuesFromPatch(string(diff), lintCtx.Log)
		if err != nil {
			return nil, fmt.Errorf("can't extract issues from gofmt diff output %q: %s", string(diff), err)
		}

		issues = append(issues, is...)
	}

	return issues, nil
}
