package goanalysis

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"go/types"
	"reflect"

	"golang.org/x/tools/go/analysis"
)

type objectFactKey struct {
	obj types.Object
	typ reflect.Type
}

type packageFactKey struct {
	pkg *types.Package
	typ reflect.Type
}

// inheritFacts populates act.facts with
// those it obtains from its dependency, dep.
func inheritFacts(act, dep *action) {
	serialize := false

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
			var err error
			fact, err = codeFact(fact)
			if err != nil {
				act.r.log.Panicf("internal error: encoding of %T fact failed in %v", fact, act)
			}
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
			var err error
			fact, err = codeFact(fact)
			if err != nil {
				act.r.log.Panicf("internal error: encoding of %T fact failed in %v", fact, act)
			}
		}

		factsInheritDebugf("%v: inherited %T fact for %s: %s", act, fact, key.pkg.Path(), fact)
		act.packageFacts[key] = fact
	}
}

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
