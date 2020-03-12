package golinters

import (
	"sort"
	"testing"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/asmdecl"
	"golang.org/x/tools/go/analysis/passes/assign"
	"golang.org/x/tools/go/analysis/passes/atomic"
	"golang.org/x/tools/go/analysis/passes/bools"
	"golang.org/x/tools/go/analysis/passes/buildtag"
	"golang.org/x/tools/go/analysis/passes/cgocall"
	"golang.org/x/tools/go/analysis/passes/shadow"

	"github.com/golangci/golangci-lint/pkg/config"
)

func TestGovet(t *testing.T) {
	// Checking that every default analyzer is in "all analyzers" list.
	checkList := append(defaultAnalyzers,
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
		if !sort.IsSorted(sortedAnalyzers(allAnalyzers)) {
			t.Error("please keep all analyzers list sorted by name")
		}
	})
	t.Run("Default", func(t *testing.T) {
		if !sort.IsSorted(sortedAnalyzers(defaultAnalyzers)) {
			t.Error("please keep default analyzers list sorted by name")
		}
	})
}

func TestGovetAnalyzerIsEnabled(t *testing.T) {
	defaultAnalyzers := []*analysis.Analyzer{
		asmdecl.Analyzer,
		assign.Analyzer,
		atomic.Analyzer,
		bools.Analyzer,
		buildtag.Analyzer,
		cgocall.Analyzer,
	}
	for _, tc := range []struct {
		Enable     []string
		Disable    []string
		EnableAll  bool
		DisableAll bool

		Name    string
		Enabled bool
	}{
		{Name: "assign", Enabled: true},
		{Name: "cgocall", Enabled: false, DisableAll: true},
		{Name: "errorsas", Enabled: false},
		{Name: "bools", Enabled: false, Disable: []string{"bools"}},
		{Name: "unsafeptr", Enabled: true, Enable: []string{"unsafeptr"}},
		{Name: "shift", Enabled: true, EnableAll: true},
		{Name: "shadow", EnableAll: true, Disable: []string{"shadow"}, Enabled: false},
	} {
		cfg := &config.GovetSettings{
			Enable:     tc.Enable,
			Disable:    tc.Disable,
			EnableAll:  tc.EnableAll,
			DisableAll: tc.DisableAll,
		}
		if enabled := isAnalyzerEnabled(tc.Name, cfg, defaultAnalyzers); enabled != tc.Enabled {
			t.Errorf("%+v", tc)
		}
	}
}
