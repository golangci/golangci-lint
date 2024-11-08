// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// Partial copy of https://github.com/golang/tools/blob/dba5486c2a1d03519930812112b23ed2c45c04fc/go/analysis/internal/checker/checker.go

package goanalysis

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"go/types"
	"reflect"
	"time"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/packages"

	"github.com/golangci/golangci-lint/pkg/goanalysis/pkgerrors"
)

// NOTE(ldez) altered: custom fields; remove 'once' and 'duration'.
// An action represents one unit of analysis work: the application of
// one analysis to one package. Actions form a DAG, both within a
// package (as different analyzers are applied, either in sequence or
// parallel), and across packages (as dependencies are analyzed).
type action struct {
	a            *analysis.Analyzer
	pkg          *packages.Package
	pass         *analysis.Pass
	isroot       bool
	deps         []*action
	objectFacts  map[objectFactKey]analysis.Fact
	packageFacts map[packageFactKey]analysis.Fact
	result       any
	diagnostics  []analysis.Diagnostic
	err          error

	// NOTE(ldez) custom fields.
	r                   *runner
	analysisDoneCh      chan struct{}
	loadCachedFactsDone bool
	loadCachedFactsOk   bool
	isInitialPkg        bool
	needAnalyzeSource   bool
}

// NOTE(ldez) no alteration.
type objectFactKey struct {
	obj types.Object
	typ reflect.Type
}

// NOTE(ldez) no alteration.
type packageFactKey struct {
	pkg *types.Package
	typ reflect.Type
}

// NOTE(ldez) no alteration.
func (act *action) String() string {
	return fmt.Sprintf("%s@%s", act.a, act.pkg)
}

// NOTE(ldez) altered version of `func (act *action) execOnce()`.
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
		if dep.err != nil {
			depErrors = errors.Join(depErrors, errors.Unwrap(dep.err))
		}
	}

	if depErrors != nil {
		act.err = fmt.Errorf("failed prerequisites: %w", depErrors)
		return
	}

	// Plumb the output values of the dependencies
	// into the inputs of this action.  Also facts.
	inputs := make(map[*analysis.Analyzer]any)
	act.objectFacts = make(map[objectFactKey]analysis.Fact)
	act.packageFacts = make(map[packageFactKey]analysis.Fact)
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

	module := &analysis.Module{} // possibly empty (non nil) in go/analysis drivers.
	if mod := act.pkg.Module; mod != nil {
		module.Path = mod.Path
		module.Version = mod.Version
		module.GoVersion = mod.GoVersion
	}

	// Run the analysis.
	pass := &analysis.Pass{
		Analyzer:     act.a,
		Fset:         act.pkg.Fset,
		Files:        act.pkg.Syntax,
		OtherFiles:   act.pkg.OtherFiles,
		IgnoredFiles: act.pkg.IgnoredFiles,
		Pkg:          act.pkg.Types,
		TypesInfo:    act.pkg.TypesInfo,
		TypesSizes:   act.pkg.TypesSizes,
		TypeErrors:   act.pkg.TypeErrors,
		Module:       module,

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

	err := act.persistFactsToCache()
	if err != nil {
		act.r.log.Warnf("Failed to persist facts to cache: %s", err)
	}
}

// NOTE(ldez) altered: logger; serialize.
// inheritFacts populates act.facts with
// those it obtains from its dependency, dep.
func inheritFacts(act, dep *action) {
	const serialize = false

	for key, fact := range dep.objectFacts {
		// Filter out facts related to objects
		// that are irrelevant downstream
		// (equivalently: not in the compiler export data).
		if !exportedFrom(key.obj, dep.pkg.Types) {
			factsInheritDebugf("%v: discarding %T fact from %s for %s: %s", act, fact, dep, key.obj, fact)
			continue
		}

		// Optionally serialize/deserialize fact
		// to verify that it works across address spaces.
		if serialize {
			encodedFact, err := codeFact(fact)
			if err != nil {
				act.r.log.Panicf("internal error: encoding of %T fact failed in %v: %v", fact, act, err)
			}
			fact = encodedFact
		}

		factsInheritDebugf("%v: inherited %T fact for %s: %s", act, fact, key.obj, fact)

		act.objectFacts[key] = fact
	}

	for key, fact := range dep.packageFacts {
		// TODO: filter out facts that belong to
		// packages not mentioned in the export data
		// to prevent side channels.

		// Optionally serialize/deserialize fact
		// to verify that it works across address spaces
		// and is deterministic.
		if serialize {
			encodedFact, err := codeFact(fact)
			if err != nil {
				act.r.log.Panicf("internal error: encoding of %T fact failed in %v", fact, act)
			}
			fact = encodedFact
		}

		factsInheritDebugf("%v: inherited %T fact for %s: %s", act, fact, key.pkg.Path(), fact)

		act.packageFacts[key] = fact
	}
}

// NOTE(ldez) no alteration.
// codeFact encodes then decodes a fact,
// just to exercise that logic.
func codeFact(fact analysis.Fact) (analysis.Fact, error) {
	// We encode facts one at a time.
	// A real modular driver would emit all facts
	// into one encoder to improve gob efficiency.
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(fact); err != nil {
		return nil, err
	}

	// Encode it twice and assert that we get the same bits.
	// This helps detect nondeterministic Gob encoding (e.g. of maps).
	var buf2 bytes.Buffer
	if err := gob.NewEncoder(&buf2).Encode(fact); err != nil {
		return nil, err
	}
	if !bytes.Equal(buf.Bytes(), buf2.Bytes()) {
		return nil, fmt.Errorf("encoding of %T fact is nondeterministic", fact)
	}

	newFact := reflect.New(reflect.TypeOf(fact).Elem()).Interface().(analysis.Fact)
	if err := gob.NewDecoder(&buf).Decode(newFact); err != nil {
		return nil, err
	}
	return newFact, nil
}

// NOTE(ldez) no alteration.
// exportedFrom reports whether obj may be visible to a package that imports pkg.
// This includes not just the exported members of pkg, but also unexported
// constants, types, fields, and methods, perhaps belonging to other packages,
// that find there way into the API.
// This is an over-approximation of the more accurate approach used by
// gc export data, which walks the type graph, but it's much simpler.
//
// TODO(adonovan): do more accurate filtering by walking the type graph.
func exportedFrom(obj types.Object, pkg *types.Package) bool {
	switch obj := obj.(type) {
	case *types.Func:
		return obj.Exported() && obj.Pkg() == pkg ||
			obj.Type().(*types.Signature).Recv() != nil
	case *types.Var:
		if obj.IsField() {
			return true
		}
		// we can't filter more aggressively than this because we need
		// to consider function parameters exported, but have no way
		// of telling apart function parameters from local variables.
		return obj.Pkg() == pkg
	case *types.TypeName, *types.Const:
		return true
	}
	return false // Nil, Builtin, Label, or PkgName
}

// NOTE(ldez) altered: logger; `act.factType`
// importObjectFact implements Pass.ImportObjectFact.
// Given a non-nil pointer ptr of type *T, where *T satisfies Fact,
// importObjectFact copies the fact value to *ptr.
func (act *action) importObjectFact(obj types.Object, ptr analysis.Fact) bool {
	if obj == nil {
		panic("nil object")
	}
	key := objectFactKey{obj, act.factType(ptr)}
	if v, ok := act.objectFacts[key]; ok {
		reflect.ValueOf(ptr).Elem().Set(reflect.ValueOf(v).Elem())
		return true
	}
	return false
}

// NOTE(ldez) altered: removes code related to `act.pass.ExportPackageFact`; logger; `act.factType`.
// exportObjectFact implements Pass.ExportObjectFact.
func (act *action) exportObjectFact(obj types.Object, fact analysis.Fact) {
	if obj.Pkg() != act.pkg.Types {
		act.r.log.Panicf("internal error: in analysis %s of package %s: Fact.Set(%s, %T): can't set facts on objects belonging another package",
			act.a, act.pkg, obj, fact)
	}

	key := objectFactKey{obj, act.factType(fact)}
	act.objectFacts[key] = fact // clobber any existing entry
	if isFactsExportDebug {
		objstr := types.ObjectString(obj, (*types.Package).Name)

		factsExportDebugf("%s: object %s has fact %s\n",
			act.pkg.Fset.Position(obj.Pos()), objstr, fact)
	}
}

// NOTE(ldez) no alteration.
func (act *action) allObjectFacts() []analysis.ObjectFact {
	facts := make([]analysis.ObjectFact, 0, len(act.objectFacts))
	for k := range act.objectFacts {
		facts = append(facts, analysis.ObjectFact{Object: k.obj, Fact: act.objectFacts[k]})
	}
	return facts
}

// NOTE(ldez) altered: `act.factType`
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

// NOTE(ldez) altered: removes code related to `act.pass.ExportPackageFact`; logger; `act.factType`.
// exportPackageFact implements Pass.ExportPackageFact.
func (act *action) exportPackageFact(fact analysis.Fact) {
	key := packageFactKey{act.pass.Pkg, act.factType(fact)}
	act.packageFacts[key] = fact // clobber any existing entry

	factsDebugf("%s: package %s has fact %s\n",
		act.pkg.Fset.Position(act.pass.Files[0].Pos()), act.pass.Pkg.Path(), fact)
}

// NOTE(ldez) altered: add receiver to handle logs.
func (act *action) factType(fact analysis.Fact) reflect.Type {
	t := reflect.TypeOf(fact)
	if t.Kind() != reflect.Ptr {
		act.r.log.Fatalf("invalid Fact type: got %T, want pointer", fact)
	}
	return t
}

// NOTE(ldez) no alteration.
func (act *action) allPackageFacts() []analysis.PackageFact {
	facts := make([]analysis.PackageFact, 0, len(act.packageFacts))
	for k := range act.packageFacts {
		facts = append(facts, analysis.PackageFact{Package: k.pkg, Fact: act.packageFacts[k]})
	}
	return facts
}
