package golinters

import (
	"sort"
	"testing"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/shadow"
)

func TestGovet(t *testing.T) {
	// Checking that every default analyzer is in "all analyzers" list.
	allAnalyzers := getAllAnalyzers()
	checkList := append(getDefaultAnalyzers(),
		shadow.Analyzer, // special case, used in analyzersFromConfig
	)
	for _, defaultAnalyzer := range checkList {
		found := false
		for _, a := range allAnalyzers {
			if a.Name == defaultAnalyzer.Name {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("%s is not in allAnalyzers", defaultAnalyzer.Name)
		}
	}
}

type sortedAnalyzers []*analysis.Analyzer

func (p sortedAnalyzers) Len() int           { return len(p) }
func (p sortedAnalyzers) Less(i, j int) bool { return p[i].Name < p[j].Name }
func (p sortedAnalyzers) Swap(i, j int)      { p[i].Name, p[j].Name = p[j].Name, p[i].Name }

func TestGovetSorted(t *testing.T) {
	// Keeping analyzers sorted so their order match the import order.
	t.Run("All", func(t *testing.T) {
		if !sort.IsSorted(sortedAnalyzers(getAllAnalyzers())) {
			t.Error("please keep all analyzers list sorted by name")
		}
	})
	t.Run("Default", func(t *testing.T) {
		if !sort.IsSorted(sortedAnalyzers(getDefaultAnalyzers())) {
			t.Error("please keep default analyzers list sorted by name")
		}
	})
}
