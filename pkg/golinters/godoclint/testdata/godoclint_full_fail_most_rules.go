//golangcitest:args -Egodoclint
//golangcitest:config_path testdata/godoclint.yml

// Asserting rule "pkg-doc"

// bad godoc // want `package godoc should start with "PACKAGE testdata "`
package testdata

// This is a special stdlib import because the package itself has issues that
// godoclint can, but must not, detect.
import "go/ast"

// Asserting rule "start-with-name"

// bad godoc // want `godoc should start with symbol name \(pattern "GODOC %"\)`
type FooType struct{}

// bad godoc // want `godoc should start with symbol name \(pattern "GODOC %"\)`
type FooAlias = ast.Comment

// bad godoc // want `godoc should start with symbol name \(pattern "GODOC %"\)`
const FooConst = 1

// bad godoc // want `godoc should start with symbol name \(pattern "GODOC %"\)`
var FooVar = 1

// bad godoc // want `godoc should start with symbol name \(pattern "GODOC %"\)`
func FooFunc() {}

// bad godoc // want `godoc should start with symbol name \(pattern "GODOC %"\)`
func (FooType) FooFunc() {}

// bad godoc // want `godoc should start with symbol name \(pattern "GODOC %"\)`
type fooType struct{}

// bad godoc // want `godoc should start with symbol name \(pattern "GODOC %"\)`
type fooAlias = ast.Comment

// bad godoc // want `godoc should start with symbol name \(pattern "GODOC %"\)`
const fooConst = 1

// bad godoc // want `godoc should start with symbol name \(pattern "GODOC %"\)`
var fooVar = 1

// bad godoc // want `godoc should start with symbol name \(pattern "GODOC %"\)`
func fooFunc() {}

// bad godoc // want `godoc should start with symbol name \(pattern "GODOC %"\)`
func (FooType) fooFunc() {}

// Asserting rule "require-doc"

// The //foo:bar directives mark the trailing comment as a directive so they're
// not parsed as a normal trailing comment group.

type BarType struct{} //foo:bar // want `symbol should have a godoc \(BarType\)`

type BarAlias = ast.Comment //foo:bar // want `symbol should have a godoc \(BarAlias\)`

const BarConst = 1 //foo:bar // want `symbol should have a godoc \(BarConst\)`

var BarVar = 1 //foo:bar // want `symbol should have a godoc \(BarVar\)`

func BarFunc() {} //foo:bar // want `symbol should have a godoc \(BarFunc\)`

func (BarType) BarFunc() {} //foo:bar // want `symbol should have a godoc \(BarFunc\)`

type barType struct{} //foo:bar // want `symbol should have a godoc \(barType\)`

type barAlias = ast.Comment //foo:bar // want `symbol should have a godoc \(barAlias\)`

const barConst = 1 //foo:bar // want `symbol should have a godoc \(barConst\)`

var barVar = 1 //foo:bar // want `symbol should have a godoc \(barVar\)`

func barFunc() {} //foo:bar // want `symbol should have a godoc \(barFunc\)`

func (BarType) barFunc() {} //foo:bar // want `symbol should have a godoc \(barFunc\)`

// Asserting rule "no-unused-link"

// GODOC constWithUnusedLink point to [used] and unused links. // want `godoc has unused link \(unused\)`
//
// [used]: https://example.com
//
// [unused]: https://example.com
const constWithUnusedLink = 1

// Asserting rule "max-len"

// GODOC constWithTooLongGodoc has a very long godoc that exceeds the maximum allowed length for godoc comments in this test setup. // want `godoc exceeds max length \(177 > 127\)`
const constWithTooLongGodoc = 1
