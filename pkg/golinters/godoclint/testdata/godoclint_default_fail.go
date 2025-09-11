//golangcitest:args -Egodoclint

// bad godoc // want `package godoc should start with "Package testdata "`
package testdata

// This is a special stdlib import because the package itself has issues that
// godoclint can, but must not, detect.
import "go/ast"

// bad godoc // want `godoc should start with symbol name \("FooType"\)`
type FooType struct{}

// bad godoc // want `godoc should start with symbol name \("FooAlias"\)`
type FooAlias = ast.Comment

// bad godoc // want `godoc should start with symbol name \("FooConst"\)`
const FooConst = 1

// bad godoc // want `godoc should start with symbol name \("FooVar"\)`
var FooVar = 1

// bad godoc // want `godoc should start with symbol name \("FooFunc"\)`
func FooFunc() {}

// bad godoc // want `godoc should start with symbol name \("FooFunc"\)`
func (FooType) FooFunc() {}

// DeprecatedConstA is... // want `deprecation note should be formatted as "Deprecated: "`
//
// DEPRECATED: do not use
const DeprecatedConstA = 1

// DeprecatedConstB is... // want `deprecation note should be formatted as "Deprecated: "`
//
// DEPRECATED:do not use
const DeprecatedConstB = 1

// DeprecatedConstC is... // want `deprecation note should be formatted as "Deprecated: "`
//
// deprecated:do not use
const DeprecatedConstC = 1
