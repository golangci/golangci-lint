package goanalysis

import (
	"bytes"
	"fmt"
	"go/token"
	"io"
	"os"

	"github.com/pkg/errors"
	"github.com/pmezard/go-difflib/difflib"

	"github.com/golangci/golangci-lint/pkg/lint/linter"
	"github.com/golangci/golangci-lint/pkg/result"
)

func convertSuggestedFixes(lintCtx *linter.Context, linterName string, diag *Diagnostic) ([]result.Issue, error) {
	var issues []result.Issue

	for _, fix := range diag.SuggestedFixes {
		origContent := make(map[string][]byte)
		readers := make(map[string]io.ReadSeeker)
		bufs := make(map[string]*bytes.Buffer)

		for _, edit := range fix.TextEdits {
			if edit.End == token.NoPos {
				edit.End = edit.Pos
			}
			start, end := diag.Pkg.Fset.Position(edit.Pos), diag.Pkg.Fset.Position(edit.End)
			orig, ok := readers[start.Filename]
			if !ok {
				data, err := os.ReadFile(start.Filename)
				if err != nil {
					return nil, errors.Wrapf(err, "can't read file: %s", start.Filename)
				}
				origContent[start.Filename] = data
				orig = bytes.NewReader(data)
				readers[start.Filename] = orig
				bufs[start.Filename] = &bytes.Buffer{}
			}

			buf := bufs[start.Filename]
			cur, err := orig.Seek(0, io.SeekCurrent)
			if err != nil {
				return nil, errors.Wrapf(err, "can't get current position of reader of file %q", start.Filename)
			}
			if l := start.Offset - int(cur); l > 0 {
				b := make([]byte, l)
				if _, err := orig.Read(b); err != nil {
					return nil, errors.Wrapf(err, "can't read form stored reader of file %q", start.Filename)
				}
				buf.Write(b)
			}
			buf.Write(edit.NewText)
			if _, err := orig.Seek(int64(end.Offset), io.SeekStart); err != nil {
				return nil, errors.Wrapf(err, "can't change position of reader of file %q", start.Filename)
			}
		}
		for filename, f := range readers {
			data, err := io.ReadAll(f)
			if err != nil {
				return nil, err
			}
			bufs[filename].Write(data)
		}

		for filename, data := range origContent {
			newIssues, err := createPatchAndExtractIssues(lintCtx, linterName, filename, data, bufs[filename].Bytes())
			if err != nil {
				return nil, err
			}
			for i := range newIssues {
				newIssues[i].Text = fmt.Sprintf("%s: %s", diag.Message, fix.Message)
			}
			issues = append(issues, newIssues...)
		}
	}

	return issues, nil
}

func createPatchAndExtractIssues(lintCtx *linter.Context, linterName, filename string, src, dst []byte) ([]result.Issue, error) {
	out := bytes.Buffer{}
	if _, err := out.WriteString(fmt.Sprintf("--- %[1]s\n+++ %[1]s\n", filename)); err != nil {
		return nil, errors.Wrap(err, "can't write diff header")
	}

	d := difflib.UnifiedDiff{
		A:       difflib.SplitLines(string(src)),
		B:       difflib.SplitLines(string(dst)),
		Context: 3,
	}

	if err := difflib.WriteUnifiedDiff(&out, d); err != nil {
		return nil, errors.Wrap(err, "can't create diff")
	}

	newIssues, err := ExtractIssuesFromPatch(out.String(), lintCtx.Log, lintCtx, linterName)
	if err != nil {
		return nil, errors.Wrap(err, "can't extract issues from diff")
	}

	return newIssues, nil
}
