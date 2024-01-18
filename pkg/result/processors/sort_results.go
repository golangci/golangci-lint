package processors

import (
	"fmt"
	"sort"
	"strings"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/result"
)

// Base propose of this functionality to sort results (issues)
// produced by various linters by analyzing code. We're achieving this
// by sorting results.Issues using processor step, and chain based
// rules that can compare different properties of the Issues struct.

var _ Processor = (*SortResults)(nil)

type SortResults struct {
	cmp comparator
	cfg *config.Config
}

func NewSortResults(cfg *config.Config) *SortResults {
	// For sorting we are comparing (in next order): file names, line numbers,
	// position, and finally - giving up.
	return &SortResults{
		cmp: ByName{
			next: ByLine{
				next: ByColumn{},
			},
		},
		cfg: cfg,
	}
}

// Process is performing sorting of the result issues.
func (sr SortResults) Process(issues []result.Issue) ([]result.Issue, error) {
	if sr.cfg.Output.GroupResultsByLinter {
		issuesByLinterName := sr.groupIssuesByLinterName(issues)
		return sr.processIssuesByLinterName(issuesByLinterName)
	}
	return sr.processFlatIssues(issues)
}

func (sr SortResults) groupIssuesByLinterName(issues []result.Issue) map[string][]result.Issue {
	issuesByLinterName := map[string][]result.Issue{}
	for _, issue := range issues {
		if _, ok := issuesByLinterName[issue.FromLinter]; !ok {
			issuesByLinterName[issue.FromLinter] = []result.Issue{}
		}
		issuesByLinterName[issue.FromLinter] = append(issuesByLinterName[issue.FromLinter], issue)
	}
	return issuesByLinterName
}

func (sr SortResults) processIssuesByLinterName(issuesByLinterName map[string][]result.Issue) ([]result.Issue, error) {
	linterNames := sr.getSortedLinterNames(issuesByLinterName)
	var processedIssues []result.Issue
	for _, linterName := range linterNames {
		linterIssues := issuesByLinterName[linterName]
		processedLinterIssues, err := sr.processFlatIssues(linterIssues)
		if err != nil {
			return nil, fmt.Errorf("failed to process issues from %s linter: %w", linterName, err)
		}
		processedIssues = append(processedIssues, processedLinterIssues...)
	}
	return processedIssues, nil
}

func (sr SortResults) getSortedLinterNames(issuesByLinterName map[string][]result.Issue) []string {
	linterNames := make([]string, 0, len(issuesByLinterName))
	for linterName := range issuesByLinterName {
		linterNames = append(linterNames, linterName)
	}
	sort.Strings(linterNames)
	return linterNames
}

func (sr SortResults) processFlatIssues(issues []result.Issue) ([]result.Issue, error) {
	if !sr.cfg.Output.SortResults {
		return issues, nil
	}

	sort.Slice(issues, func(i, j int) bool {
		return sr.cmp.Compare(&issues[i], &issues[j]) == Less
	})

	return issues, nil
}

func (sr SortResults) Name() string { return "sort_results" }
func (sr SortResults) Finish()      {}

type compareResult int

const (
	Less compareResult = iota - 1
	Equal
	Greater
	None
)

func (c compareResult) isNeutral() bool {
	// return true if compare result is incomparable or equal.
	return c == None || c == Equal
}

func (c compareResult) String() string {
	switch c {
	case Less:
		return "Less"
	case Equal:
		return "Equal"
	case Greater:
		return "Greater"
	}

	return "None"
}

// comparator describe how to implement compare for two "issues" lexicographically
type comparator interface {
	Compare(a, b *result.Issue) compareResult
	Next() comparator
}

var (
	_ comparator = (*ByName)(nil)
	_ comparator = (*ByLine)(nil)
	_ comparator = (*ByColumn)(nil)
)

type ByName struct{ next comparator }

func (cmp ByName) Next() comparator { return cmp.next }

func (cmp ByName) Compare(a, b *result.Issue) compareResult {
	var res compareResult

	if res = compareResult(strings.Compare(a.FilePath(), b.FilePath())); !res.isNeutral() {
		return res
	}

	if next := cmp.Next(); next != nil {
		return next.Compare(a, b)
	}

	return res
}

type ByLine struct{ next comparator }

func (cmp ByLine) Next() comparator { return cmp.next }

func (cmp ByLine) Compare(a, b *result.Issue) compareResult {
	var res compareResult

	if res = numericCompare(a.Line(), b.Line()); !res.isNeutral() {
		return res
	}

	if next := cmp.Next(); next != nil {
		return next.Compare(a, b)
	}

	return res
}

type ByColumn struct{ next comparator }

func (cmp ByColumn) Next() comparator { return cmp.next }

func (cmp ByColumn) Compare(a, b *result.Issue) compareResult {
	var res compareResult

	if res = numericCompare(a.Column(), b.Column()); !res.isNeutral() {
		return res
	}

	if next := cmp.Next(); next != nil {
		return next.Compare(a, b)
	}

	return res
}

func numericCompare(a, b int) compareResult {
	var (
		isValuesInvalid  = a < 0 || b < 0
		isZeroValuesBoth = a == 0 && b == 0
		isEqual          = a == b
		isZeroValueInA   = b > 0 && a == 0
		isZeroValueInB   = a > 0 && b == 0
	)

	switch {
	case isZeroValuesBoth || isEqual:
		return Equal
	case isValuesInvalid || isZeroValueInA || isZeroValueInB:
		return None
	case a > b:
		return Greater
	case a < b:
		return Less
	}

	return Equal
}
