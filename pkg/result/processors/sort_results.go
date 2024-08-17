package processors

import (
	"errors"
	"fmt"
	"slices"
	"sort"
	"strings"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/result"
)

// Base propose of this functionality to sort results (issues)
// produced by various linters by analyzing code. We're achieving this
// by sorting results.Issues using processor step, and chain based
// rules that can compare different properties of the Issues struct.

const (
	orderNameFile     = "file"
	orderNameLinter   = "linter"
	orderNameSeverity = "severity"
)

var _ Processor = (*SortResults)(nil)

type SortResults struct {
	cmps map[string]*comparator

	cfg *config.Output
}

func NewSortResults(cfg *config.Config) *SortResults {
	return &SortResults{
		cmps: map[string]*comparator{
			// For sorting we are comparing (in next order):
			// file names, line numbers, position, and finally - giving up.
			orderNameFile: byFileName().SetNext(byLine().SetNext(byColumn())),
			// For sorting we are comparing: linter name
			orderNameLinter: byLinter(),
			// For sorting we are comparing: severity
			orderNameSeverity: bySeverity(),
		},
		cfg: &cfg.Output,
	}
}

func (SortResults) Name() string { return "sort_results" }

// Process is performing sorting of the result issues.
func (p SortResults) Process(issues []result.Issue) ([]result.Issue, error) {
	if !p.cfg.SortResults {
		return issues, nil
	}

	if len(p.cfg.SortOrder) == 0 {
		p.cfg.SortOrder = []string{orderNameFile}
	}

	var cmps []*comparator
	for _, name := range p.cfg.SortOrder {
		c, ok := p.cmps[name]
		if !ok {
			return nil, fmt.Errorf("unsupported sort-order name %q", name)
		}

		cmps = append(cmps, c)
	}

	cmp, err := mergeComparators(cmps)
	if err != nil {
		return nil, err
	}

	sort.Slice(issues, func(i, j int) bool {
		return cmp.Compare(&issues[i], &issues[j]) == less
	})

	return issues, nil
}

func (SortResults) Finish() {}

type compareResult int

const (
	less compareResult = iota - 1
	equal
	greater
	none
)

func (c compareResult) isNeutral() bool {
	// return true if compare result is incomparable or equal.
	return c == none || c == equal
}

func (c compareResult) String() string {
	switch c {
	case less:
		return "less"
	case equal:
		return "equal"
	case greater:
		return "greater"
	default:
		return "none"
	}
}

// comparator describes how to implement compare for two "issues".
type comparator struct {
	name    string
	compare func(a, b *result.Issue) compareResult
	next    *comparator
}

func (cmp *comparator) Next() *comparator { return cmp.next }

func (cmp *comparator) SetNext(c *comparator) *comparator {
	cmp.next = c
	return cmp
}

func (cmp *comparator) String() string {
	s := cmp.name
	if cmp.Next() != nil {
		s += " > " + cmp.Next().String()
	}

	return s
}

func (cmp *comparator) Compare(a, b *result.Issue) compareResult {
	res := cmp.compare(a, b)
	if !res.isNeutral() {
		return res
	}

	if next := cmp.Next(); next != nil {
		return next.Compare(a, b)
	}

	return res
}

func byFileName() *comparator {
	return &comparator{
		name: "byFileName",
		compare: func(a, b *result.Issue) compareResult {
			return compareResult(strings.Compare(a.FilePath(), b.FilePath()))
		},
	}
}

func byLine() *comparator {
	return &comparator{
		name: "byLine",
		compare: func(a, b *result.Issue) compareResult {
			return numericCompare(a.Line(), b.Line())
		},
	}
}

func byColumn() *comparator {
	return &comparator{
		name: "byColumn",
		compare: func(a, b *result.Issue) compareResult {
			return numericCompare(a.Column(), b.Column())
		},
	}
}

func byLinter() *comparator {
	return &comparator{
		name: "byLinter",
		compare: func(a, b *result.Issue) compareResult {
			return compareResult(strings.Compare(a.FromLinter, b.FromLinter))
		},
	}
}

func bySeverity() *comparator {
	return &comparator{
		name: "bySeverity",
		compare: func(a, b *result.Issue) compareResult {
			return severityCompare(a.Severity, b.Severity)
		},
	}
}

func mergeComparators(cmps []*comparator) (*comparator, error) {
	if len(cmps) == 0 {
		return nil, errors.New("no comparator")
	}

	for i := range len(cmps) - 1 {
		findComparatorTip(cmps[i]).SetNext(cmps[i+1])
	}

	return cmps[0], nil
}

func findComparatorTip(cmp *comparator) *comparator {
	if cmp.Next() != nil {
		return findComparatorTip(cmp.Next())
	}

	return cmp
}

func severityCompare(a, b string) compareResult {
	// The position inside the slice define the importance (lower to higher).
	classic := []string{"low", "medium", "high", "warning", "error"}

	if slices.Contains(classic, a) && slices.Contains(classic, b) {
		switch {
		case slices.Index(classic, a) > slices.Index(classic, b):
			return greater
		case slices.Index(classic, a) < slices.Index(classic, b):
			return less
		default:
			return equal
		}
	}

	if slices.Contains(classic, a) {
		return greater
	}

	if slices.Contains(classic, b) {
		return less
	}

	return compareResult(strings.Compare(a, b))
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
		return equal
	case isValuesInvalid || isZeroValueInA || isZeroValueInB:
		return none
	case a > b:
		return greater
	case a < b:
		return less
	}

	return equal
}
