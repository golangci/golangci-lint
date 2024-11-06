package goanalysis

import (
	"fmt"
	"go/types"
	"reflect"
	"runtime/debug"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/packages"

	"github.com/golangci/golangci-lint/internal/errorutil"
)

type actionAllocator struct {
	allocatedActions []action
	nextFreeIndex    int
}

func newActionAllocator(maxCount int) *actionAllocator {
	return &actionAllocator{
		allocatedActions: make([]action, maxCount),
		nextFreeIndex:    0,
	}
}

func (actAlloc *actionAllocator) alloc() *action {
	if actAlloc.nextFreeIndex == len(actAlloc.allocatedActions) {
		panic(fmt.Sprintf("Made too many allocations of actions: %d allowed", len(actAlloc.allocatedActions)))
	}
	act := &actAlloc.allocatedActions[actAlloc.nextFreeIndex]
	actAlloc.nextFreeIndex++
	return act
}

// An action represents one unit of analysis work: the application of
// one analysis to one package. Actions form a DAG, both within a
// package (as different analyzers are applied, either in sequence or
// parallel), and across packages (as dependencies are analyzed).
type action struct {
	a                   *analysis.Analyzer
	pkg                 *packages.Package
	pass                *analysis.Pass
	deps                []*action
	objectFacts         map[objectFactKey]analysis.Fact
	packageFacts        map[packageFactKey]analysis.Fact
	result              any
	diagnostics         []analysis.Diagnostic
	err                 error
	r                   *runner
	analysisDoneCh      chan struct{}
	loadCachedFactsDone bool
	loadCachedFactsOk   bool
	isroot              bool
	isInitialPkg        bool
	needAnalyzeSource   bool
}

func (act *action) waitUntilDependingAnalyzersWorked() {
	for _, dep := range act.deps {
		if dep.pkg == act.pkg {
			<-dep.analysisDoneCh
		}
	}
}

func (act *action) analyzeSafe() {
	defer func() {
		if p := recover(); p != nil {
			if !act.isroot {
				// This line allows to display "hidden" panic with analyzers like buildssa.
				// Some linters are dependent of sub-analyzers but when a sub-analyzer fails the linter is not aware of that,
				// this results to another panic (ex: "interface conversion: interface {} is nil, not *buildssa.SSA").
				act.r.log.Errorf("%s: panic during analysis: %v, %s", act.a.Name, p, string(debug.Stack()))
			}

			act.err = errorutil.NewPanicError(fmt.Sprintf("%s: package %q (isInitialPkg: %t, needAnalyzeSource: %t): %s",
				act.a.Name, act.pkg.Name, act.isInitialPkg, act.needAnalyzeSource, p), debug.Stack())
		}
	}()

	act.r.sw.TrackStage(act.a.Name, act.analyze)
}

// importPackageFact implements Pass.ImportPackageFact.
// Given a non-nil pointer ptr of type *T, where *T satisfies Fact,
// fact copies the fact value to *ptr.
func (act *action) importPackageFact(pkg *types.Package, ptr analysis.Fact) bool {
	if pkg == nil {
		panic("nil package")
	}
	key := packageFactKey{pkg, act.factType(ptr)}
	if v, ok := act.packageFacts[key]; ok {
		reflect.ValueOf(ptr).Elem().Set(reflect.ValueOf(v).Elem())
		return true
	}
	return false
}

func (act *action) markDepsForAnalyzingSource() {
	// Horizontal deps (analyzer.Requires) must be loaded from source and analyzed before analyzing
	// this action.
	for _, dep := range act.deps {
		if dep.pkg == act.pkg {
			// Analyze source only for horizontal dependencies, e.g. from "buildssa".
			dep.needAnalyzeSource = true // can't be set in parallel
		}
	}
}
