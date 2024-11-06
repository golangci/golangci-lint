package goanalysis

import (
	"errors"
	"fmt"
	"go/types"
	"reflect"
	"runtime/debug"
	"time"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/packages"

	"github.com/golangci/golangci-lint/internal/errorutil"
	"github.com/golangci/golangci-lint/pkg/goanalysis/pkgerrors"
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

func (act *action) String() string {
	return fmt.Sprintf("%s@%s", act.a, act.pkg)
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

func (act *action) analyze() {
	defer close(act.analysisDoneCh) // unblock actions depending on this action

	if !act.needAnalyzeSource {
		return
	}

	defer func(now time.Time) {
		analyzeDebugf("go/analysis: %s: %s: analyzed package %q in %s", act.r.prefix, act.a.Name, act.pkg.Name, time.Since(now))
	}(time.Now())

	// Report an error if any dependency failures.
	var depErrors error
	for _, dep := range act.deps {
		if dep.err == nil {
			continue
		}

		depErrors = errors.Join(depErrors, errors.Unwrap(dep.err))
	}
	if depErrors != nil {
		act.err = fmt.Errorf("failed prerequisites: %w", depErrors)
		return
	}

	// Plumb the output values of the dependencies
	// into the inputs of this action.  Also facts.
	inputs := make(map[*analysis.Analyzer]any)
	startedAt := time.Now()
	for _, dep := range act.deps {
		if dep.pkg == act.pkg {
			// Same package, different analysis (horizontal edge):
			// in-memory outputs of prerequisite analyzers
			// become inputs to this analysis pass.
			inputs[dep.a] = dep.result
		} else if dep.a == act.a { // (always true)
			// Same analysis, different package (vertical edge):
			// serialized facts produced by prerequisite analysis
			// become available to this analysis pass.
			inheritFacts(act, dep)
		}
	}
	factsDebugf("%s: Inherited facts in %s", act, time.Since(startedAt))

	// Run the analysis.
	pass := &analysis.Pass{
		Analyzer:          act.a,
		Fset:              act.pkg.Fset,
		Files:             act.pkg.Syntax,
		OtherFiles:        act.pkg.OtherFiles,
		Pkg:               act.pkg.Types,
		TypesInfo:         act.pkg.TypesInfo,
		TypesSizes:        act.pkg.TypesSizes,
		ResultOf:          inputs,
		Report:            func(d analysis.Diagnostic) { act.diagnostics = append(act.diagnostics, d) },
		ImportObjectFact:  act.importObjectFact,
		ExportObjectFact:  act.exportObjectFact,
		ImportPackageFact: act.importPackageFact,
		ExportPackageFact: act.exportPackageFact,
		AllObjectFacts:    act.allObjectFacts,
		AllPackageFacts:   act.allPackageFacts,
	}
	act.pass = pass
	act.r.passToPkgGuard.Lock()
	act.r.passToPkg[pass] = act.pkg
	act.r.passToPkgGuard.Unlock()

	if act.pkg.IllTyped {
		// It looks like there should be !pass.Analyzer.RunDespiteErrors
		// but govet's cgocall crashes on it. Govet itself contains !pass.Analyzer.RunDespiteErrors condition here,
		// but it exits before it if packages.Load have failed.
		act.err = fmt.Errorf("analysis skipped: %w", &pkgerrors.IllTypedError{Pkg: act.pkg})
	} else {
		startedAt = time.Now()
		act.result, act.err = pass.Analyzer.Run(pass)
		analyzedIn := time.Since(startedAt)
		if analyzedIn > time.Millisecond*10 {
			debugf("%s: run analyzer in %s", act, analyzedIn)
		}
	}

	// disallow calls after Run
	pass.ExportObjectFact = nil
	pass.ExportPackageFact = nil

	if err := act.persistFactsToCache(); err != nil {
		act.r.log.Warnf("Failed to persist facts to cache: %s", err)
	}
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
