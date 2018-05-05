package processors

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/bradleyfalzon/revgrep"
	"github.com/golangci/golangci-lint/pkg/result"
)

type DiffProcessor struct {
	patch string
}

func NewDiffProcessor(patch string) *DiffProcessor {
	return &DiffProcessor{
		patch: patch,
	}
}

func (p DiffProcessor) Name() string {
	return "diff"
}

func (p DiffProcessor) processResult(res result.Result) (*result.Result, error) {
	// Make mapping to restore original issues metadata later
	fli := makeFilesToLinesToIssuesMap([]result.Result{res})

	rIssues, err := p.runRevgrepOnIssues(res.Issues)
	if err != nil {
		return nil, err
	}

	newIssues := []result.Issue{}
	for _, ri := range rIssues {
		if fli[ri.File] == nil {
			return nil, fmt.Errorf("can't get original issue file for %v", ri)
		}

		oi := fli[ri.File][ri.LineNo]
		if len(oi) != 1 {
			return nil, fmt.Errorf("can't get original issue for %v: %v", ri, oi)
		}

		i := result.Issue{
			File:       ri.File,
			LineNumber: ri.LineNo,
			Text:       ri.Message,
			HunkPos:    ri.HunkPos,
			FromLinter: oi[0].FromLinter,
		}
		newIssues = append(newIssues, i)
	}

	res.Issues = newIssues
	return &res, nil
}

func (p DiffProcessor) Process(results []result.Result) ([]result.Result, error) {
	retResults := []result.Result{}
	for _, res := range results {
		newRes, err := p.processResult(res)
		if err != nil {
			return nil, fmt.Errorf("can't filter only new issues for result %+v: %s", res, err)
		}
		retResults = append(retResults, *newRes)
	}

	return retResults, nil
}

func (p DiffProcessor) runRevgrepOnIssues(issues []result.Issue) ([]revgrep.Issue, error) {
	// TODO: change revgrep to accept interface with line number, file name
	fakeIssuesLines := []string{}
	for _, i := range issues {
		line := fmt.Sprintf("%s:%d:%d: %s", i.File, i.LineNumber, 0, i.Text)
		fakeIssuesLines = append(fakeIssuesLines, line)
	}
	fakeIssuesOut := strings.Join(fakeIssuesLines, "\n")

	checker := revgrep.Checker{
		Patch:  strings.NewReader(p.patch),
		Regexp: `^([^:]+):(\d+):(\d+)?:?\s*(.*)$`,
	}
	rIssues, err := checker.Check(strings.NewReader(fakeIssuesOut), ioutil.Discard)
	if err != nil {
		return nil, fmt.Errorf("can't filter only new issues by revgrep: %s", err)
	}

	return rIssues, nil
}
