//golangcitest:args -Egodoclint
//golangcitest:expected_exitcode 0

// Package testdata has good godoc.
package testdata

// This is a special stdlib import because the package itself has issues that
// godoclint can, but must not, detect.
import "go/ast"

// FooType is a type.
type FooType struct{}

// FooAlias is an alias.
type FooAlias = ast.Comment

// FooConst is a constant.
const FooConst = 1

// FooVar is a variable.
var FooVar = 1

// FooFunc is a function.
func FooFunc() {}

// FooFunc is a method.
func (FooType) FooFunc() {}

// bad godoc on unexported symbol
type fooType struct{}

// bad godoc on unexported symbol
type fooAlias = ast.Comment

// bad godoc on unexported symbol
const fooConst = 1

// bad godoc on unexported symbol
var fooVar = 1

// bad godoc on unexported symbol
func fooFunc() {}

// bad godoc on unexported symbol
func (FooType) fooFunc() {}

// DeprecatedConstA is...
//
// Deprecated: do not use
const DeprecatedConstA = 1

// Deprecated: do not use
const DeprecatedConstB = 1

// deprecatedConstC is...
//
// DEPRECATED: invalid deprecation note but okay since the symbol is not exported
const deprecatedConstC = 1
