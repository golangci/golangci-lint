package golinters

import (
	"bytes"
	"context"
	"fmt"

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

func getFirstDeletedLineNumberInHunk(h *diff.Hunk) (int, error) {
	lines := bytes.Split(h.Body, []byte{'\n'})
	lineNumber := int(h.OrigStartLine - 1)
	for _, line := range lines {
		lineNumber++

		if len(line) == 0 {
			continue
		}
		if line[0] == '-' {
			return lineNumber, nil
		}
	}

	return 0, fmt.Errorf("didn't find deletion line in hunk %s", string(h.Body))
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
			lineNumber, err := getFirstDeletedLineNumberInHunk(hunk)
			if err != nil {
				analytics.Log(context.TODO()).Infof("Can't get first deleted line number for hunk: %s", err)
				lineNumber = int(hunk.OrigStartLine) // use first line if no deletions:
			}

			text := "File is not gofmt-ed with -s"
			if g.useGoimports {
				text = "File is not goimports-ed"
			}
			i := result.Issue{
				FromLinter: g.Name(),
				File:       d.NewName,
				LineNumber: lineNumber,
				Text:       text,
			}
			issues = append(issues, i)
		}
	}

	return issues, nil
}

func (g gofmt) Run(ctx context.Context, exec executors.Executor) (*result.Result, error) {
	paths, err := getPathsForGoProject(exec.WorkDir())
	if err != nil {
		return nil, fmt.Errorf("can't get files to analyze: %s", err)
	}

	args := []string{"-d"}
	if !g.useGoimports {
		args = append(args, "-s")
	}
	args = append(args, paths.files...)
	out, err := exec.Run(ctx, g.Name(), args...)
	if err != nil {
		return nil, fmt.Errorf("can't run gofmt: %s, %s", err, out)
	}

	if len(out) == 0 { // no diff => no issues
		return &result.Result{
			Issues: []result.Issue{},
		}, nil
	}

	issues, err := g.extractIssuesFromPatch(out)
	if err != nil {
		return nil, fmt.Errorf("can't extract issues from gofmt diff output %q: %s", out, err)
	}

	return &result.Result{
		Issues:           issues,
		MaxIssuesPerFile: 1, // don't disturb user: show just first changed not gofmt-ed line
	}, nil
}
