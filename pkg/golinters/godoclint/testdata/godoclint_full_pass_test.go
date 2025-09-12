//golangcitest:args -Egodoclint
//golangcitest:config_path testdata/godoclint.yml
//golangcitest:expected_exitcode 0

// Asserting rule "pkg-doc" and "require-pkg-doc" since the package has a godoc.

// Package testdata
package testdata

// This is a special stdlib import because the package itself has issues that
// godoclint can, but must not, detect.
import "go/ast"

// Asserting rule "start-with-name" and "require-doc" (since all have godocs)

// FooType is...
type FooType struct{}

// FooAlias is...
type FooAlias = ast.Comment

// FooConst is...
const FooConst = 1

// FooVar is...
var FooVar = 1

// FooFunc is...
func FooFunc() {}

// FooFunc is...
func (FooType) FooFunc() {}

// fooType is...
type fooType struct{}

// fooAlias is...
type fooAlias = ast.Comment

// fooConst is...
const fooConst = 1

// fooVar is...
var fooVar = 1

// fooFunc is...
func fooFunc() {}

// fooFunc is...
func (FooType) fooFunc() {}

// Asserting rule "no-unused-link"

// constWithUnusedLink point to a [used] link and has no unused one.
//
// [used]: https://example.com
const constWithUnusedLink = 1

// Asserting rule "max-len"

// constWithTooLongGodoc has a very long godoc that does not exceed the maximum allowed length for godoc comments.
const constWithTooLongGodoc = 1

// DeprecatedConstA is...
//
// Deprecated: do not use
const DeprecatedConstA = 1

// Deprecated: do not use
const DeprecatedConstB = 1
