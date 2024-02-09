package golinters

import (
	"slices"
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
	checkList := append([]*analysis.Analyzer{}, defaultAnalyzers...)
	checkList = append(checkList, shadow.Analyzer) // special case, used in analyzersFromConfig

	for _, defaultAnalyzer := range checkList {
		found := slices.ContainsFunc(allAnalyzers, func(a *analysis.Analyzer) bool {
			return a.Name == defaultAnalyzer.Name
		})
		if !found {
			t.Errorf("%s is not in allAnalyzers", defaultAnalyzer.Name)
		}
	}
}

func sortAnalyzers(a, b *analysis.Analyzer) int {
	if a.Name < b.Name {
		return -1
	}

	if a.Name > b.Name {
		return 1
	}

	return 0
}

func TestGovetSorted(t *testing.T) {
	// Keeping analyzers sorted so their order match the import order.
	t.Run("All", func(t *testing.T) {
		if !slices.IsSortedFunc(allAnalyzers, sortAnalyzers) {
			t.Error("please keep all analyzers list sorted by name")
		}
	})

	t.Run("Default", func(t *testing.T) {
		if !slices.IsSortedFunc(defaultAnalyzers, sortAnalyzers) {
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
		Go         string

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
		{Name: "loopclosure", EnableAll: true, Enabled: false, Go: "1.22"}, // TODO(ldez) remove loopclosure when go1.23
	} {
		cfg := &config.GovetSettings{
			Enable:     tc.Enable,
			Disable:    tc.Disable,
			EnableAll:  tc.EnableAll,
			DisableAll: tc.DisableAll,
			Go:         tc.Go,
		}
		if enabled := isAnalyzerEnabled(tc.Name, cfg, defaultAnalyzers); enabled != tc.Enabled {
			t.Errorf("%+v", tc)
		}
	}
}
