//golangcitest:args -Egodoclint
//golangcitest:config_path testdata/godoclint.yml
//golangcitest:expected_exitcode 0

// Asserting rule "pkg-doc"

// PACKAGE testdata
package testdata

// This is a special stdlib import because the package itself has issues that
// godoclint can, but must not, detect.
import "go/ast"

// Asserting rule "start-with-name" (also covering "require-doc" since all have godocs)

// GODOC FooType is...
type FooType struct{}

// GODOC FooAlias is...
type FooAlias = ast.Comment

// GODOC FooConst is...
const FooConst = 1

// GODOC FooVar is...
var FooVar = 1

// GODOC FooFunc is...
func FooFunc() {}

// GODOC FooFunc is...
func (FooType) FooFunc() {}

// GODOC fooType is...
type fooType struct{}

// GODOC fooAlias is...
type fooAlias = ast.Comment

// GODOC fooConst is...
const fooConst = 1

// GODOC fooVar is...
var fooVar = 1

// GODOC fooFunc is...
func fooFunc() {}

// GODOC fooFunc is...
func (FooType) fooFunc() {}

// Asserting rule "no-unused-link"

// GODOC constWithUnusedLink point to a [used] link and has no unused one.
//
// [used]: https://example.com
const constWithUnusedLink = 1

// Asserting rule "max-len"

// GODOC constWithTooLongGodoc has a very long godoc that does not exceed the maximum allowed length for godoc comments.
const constWithTooLongGodoc = 1
