// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// Partial copy of https://github.com/golang/tools/blob/master/go/analysis/internal/checker
// FIXME add a commit hash.

package goanalysis

import (
	"reflect"

	"golang.org/x/tools/go/analysis"
)

// NOTE(ldez) no alteration.
func (act *action) allObjectFacts() []analysis.ObjectFact {
	out := make([]analysis.ObjectFact, 0, len(act.objectFacts))
	for key, fact := range act.objectFacts {
		out = append(out, analysis.ObjectFact{
			Object: key.obj,
			Fact:   fact,
		})
	}
	return out
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
		act.r.log.Fatalf("invalid Fact type: got %T, want pointer", t)
	}
	return t
}

// NOTE(ldez) no alteration.
func (act *action) allPackageFacts() []analysis.PackageFact {
	out := make([]analysis.PackageFact, 0, len(act.packageFacts))
	for key, fact := range act.packageFacts {
		out = append(out, analysis.PackageFact{
			Package: key.pkg,
			Fact:    fact,
		})
	}
	return out
}
