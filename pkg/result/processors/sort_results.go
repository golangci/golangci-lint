package processors

import (
	"cmp"
	"errors"
	"fmt"
	"slices"
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

const (
	less = iota - 1
	equal
	greater
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

	comp, err := mergeComparators(cmps)
	if err != nil {
		return nil, err
	}

	slices.SortFunc(issues, func(a, b result.Issue) int {
		return comp.Compare(&a, &b)
	})

	return issues, nil
}

func (SortResults) Finish() {}

// comparator describes how to implement compare for two "issues".
type comparator struct {
	name    string
	compare func(a, b *result.Issue) int
	next    *comparator
}

func (cp *comparator) Next() *comparator { return cp.next }

func (cp *comparator) SetNext(c *comparator) *comparator {
	cp.next = c
	return cp
}

func (cp *comparator) String() string {
	s := cp.name
	if cp.Next() != nil {
		s += " > " + cp.Next().String()
	}

	return s
}

func (cp *comparator) Compare(a, b *result.Issue) int {
	res := cp.compare(a, b)
	if res != equal {
		return res
	}

	if next := cp.Next(); next != nil {
		return next.Compare(a, b)
	}

	return res
}

func byFileName() *comparator {
	return &comparator{
		name: "byFileName",
		compare: func(a, b *result.Issue) int {
			return strings.Compare(a.FilePath(), b.FilePath())
		},
	}
}

func byLine() *comparator {
	return &comparator{
		name: "byLine",
		compare: func(a, b *result.Issue) int {
			return numericCompare(a.Line(), b.Line())
		},
	}
}

func byColumn() *comparator {
	return &comparator{
		name: "byColumn",
		compare: func(a, b *result.Issue) int {
			return numericCompare(a.Column(), b.Column())
		},
	}
}

func byLinter() *comparator {
	return &comparator{
		name: "byLinter",
		compare: func(a, b *result.Issue) int {
			return strings.Compare(a.FromLinter, b.FromLinter)
		},
	}
}

func bySeverity() *comparator {
	return &comparator{
		name: "bySeverity",
		compare: func(a, b *result.Issue) int {
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

func severityCompare(a, b string) int {
	// The position inside the slice define the importance (lower to higher).
	classic := []string{"low", "medium", "high", "warning", "error"}

	if slices.Contains(classic, a) && slices.Contains(classic, b) {
		return cmp.Compare(slices.Index(classic, a), slices.Index(classic, b))
	}

	if slices.Contains(classic, a) {
		return greater
	}

	if slices.Contains(classic, b) {
		return less
	}

	return strings.Compare(a, b)
}

func numericCompare(a, b int) int {
	// Negative value and 0 are skipped because they either "neutral" (default value) or "invalid.
	if a <= 0 || b <= 0 {
		return equal
	}

	return cmp.Compare(a, b)
}
