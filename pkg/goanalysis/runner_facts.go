package goanalysis

import (
	"go/types"
	"reflect"
)

type objectFactKey struct {
	obj types.Object
	typ reflect.Type
}

type packageFactKey struct {
	pkg *types.Package
	typ reflect.Type
}
