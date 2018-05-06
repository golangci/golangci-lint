package golinters

import (
	"bytes"
	"context"
	"fmt"

	gofmtAPI "github.com/golangci/gofmt/gofmt"
	goimportsAPI "github.com/golangci/gofmt/goimports"
	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/result"
	"github.com/golangci/golangci-shared/pkg/analytics"
	"github.com/golangci/golangci-shared/pkg/executors"
	"sourcegraph.com/sourcegraph/go-diff/diff"
)

type gofmt struct {
	useGoimports bool
}

func (g gofmt) Name() string {
	if g.useGoimports {
		return "goimports"
	}

	return "gofmt"
}

func getFirstDeletedAndAddedLineNumberInHunk(h *diff.Hunk) (int, int, error) {
	lines := bytes.Split(h.Body, []byte{'\n'})
	lineNumber := int(h.OrigStartLine - 1)
	firstAddedLineNumber := -1
	for _, line := range lines {
		lineNumber++

		if len(line) == 0 {
			continue
		}
		if line[0] == '+' && firstAddedLineNumber == -1 {
			firstAddedLineNumber = lineNumber
		}
		if line[0] == '-' {
			return lineNumber, firstAddedLineNumber, nil
		}
	}

	return 0, firstAddedLineNumber, fmt.Errorf("didn't find deletion line in hunk %s", string(h.Body))
}

func (g gofmt) extractIssuesFromPatch(patch string) ([]result.Issue, error) {
	diffs, err := diff.ParseMultiFileDiff([]byte(patch))
	if err != nil {
		return nil, fmt.Errorf("can't parse patch: %s", err)
	}

	if len(diffs) == 0 {
		return nil, fmt.Errorf("got no diffs from patch parser: %v", diffs)
	}

	issues := []result.Issue{}
	for _, d := range diffs {
		if len(d.Hunks) == 0 {
			analytics.Log(context.TODO()).Warnf("Got no hunks in diff %+v", d)
			continue
		}

		for _, hunk := range d.Hunks {
			deletedLine, addedLine, err := getFirstDeletedAndAddedLineNumberInHunk(hunk)
			if err != nil {
				analytics.Log(context.TODO()).Infof("Can't get first deleted line number for hunk: %s", err)
				if addedLine > 1 {
					deletedLine = addedLine - 1 // use previous line, TODO: use both prev and next lines
				} else {
					deletedLine = 1
				}
			}

			text := "File is not gofmt-ed with -s"
			if g.useGoimports {
				text = "File is not goimports-ed"
			}
			i := result.Issue{
				FromLinter: g.Name(),
				File:       d.NewName,
				LineNumber: deletedLine,
				Text:       text,
			}
			issues = append(issues, i)
		}
	}

	return issues, nil
}

func (g gofmt) Run(ctx context.Context, exec executors.Executor, cfg *config.Run) (*result.Result, error) {
	var issues []result.Issue

	for _, f := range cfg.Paths.Files {
		var diff []byte
		var err error
		if g.useGoimports {
			diff, err = goimportsAPI.Run(f)
		} else {
			diff, err = gofmtAPI.Run(f, cfg.Gofmt.Simplify)
		}
		if err != nil { // TODO: skip
			return nil, err
		}
		if diff == nil {
			continue
		}

		is, err := g.extractIssuesFromPatch(string(diff))
		if err != nil {
			return nil, fmt.Errorf("can't extract issues from gofmt diff output %q: %s", string(diff), err)
		}

		issues = append(issues, is...)
	}

	return &result.Result{
		Issues:           issues,
		MaxIssuesPerFile: 1, // don't disturb user: show just first changed not gofmt-ed line
	}, nil
}
