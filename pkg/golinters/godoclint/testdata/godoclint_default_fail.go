//golangcitest:args -Egodoclint

// bad godoc // want `package godoc should start with "Package testdata "`
package testdata

// This is a special stdlib import because the package itself has issues that
// godoclint can, but must not, detect.
import "go/ast"

// bad godoc // want `godoc should start with symbol name \(pattern "\(\(A\|a\|An\|an\|THE\|The\|the\) \)\?%"\)`
type FooType struct{}

// bad godoc // want `godoc should start with symbol name \(pattern "\(\(A\|a\|An\|an\|THE\|The\|the\) \)\?%"\)`
type FooAlias = ast.Comment

// bad godoc // want `godoc should start with symbol name \(pattern "\(\(A\|a\|An\|an\|THE\|The\|the\) \)\?%"\)`
const FooConst = 1

// bad godoc // want `godoc should start with symbol name \(pattern "\(\(A\|a\|An\|an\|THE\|The\|the\) \)\?%"\)`
var FooVar = 1

// bad godoc // want `godoc should start with symbol name \(pattern "\(\(A\|a\|An\|an\|THE\|The\|the\) \)\?%"\)`
func FooFunc() {}

// bad godoc // want `godoc should start with symbol name \(pattern "\(\(A\|a\|An\|an\|THE\|The\|the\) \)\?%"\)`
func (FooType) FooFunc() {}
