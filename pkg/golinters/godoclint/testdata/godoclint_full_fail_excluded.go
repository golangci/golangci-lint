//golangcitest:args -Egodoclint
//golangcitest:config_path testdata/godoclint.yml
//golangcitest:expected_exitcode 0

// Since this file is excluded in the config, godoclint should not report any
// issues (i.e.exit code 0).

// bad godoc
package testdata

// This is a special stdlib import because the package itself has issues that
// godoclint can, but must not, detect.
import "go/ast"

// bad godoc
type FooAlias = ast.Comment
