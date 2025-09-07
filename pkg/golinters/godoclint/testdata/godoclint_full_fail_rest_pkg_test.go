//golangcitest:args -Egodoclint
//golangcitest:config_path testdata/godoclint.yml

package testdata_test // want `package should have a godoc`

// This is a special stdlib import because the package itself has issues that
// godoclint can, but must not, detect.
import _ "go/ast"
