//golangcitest:args -Egolint --internal-cmd-test
package testdata

var Go_lint string // want "don't use underscores in Go names; var `Go_lint` should be `GoLint`"

func ExportedFuncWithNoComment() {
}

var ExportedVarWithNoComment string

type ExportedStructWithNoComment struct{}

type ExportedInterfaceWithNoComment interface{}

// Bad comment
func ExportedFuncWithBadComment() {}

type GolintTest struct{}

func (receiver1 GolintTest) A() {}

func (receiver2 GolintTest) B() {} // want "receiver name receiver2 should be consistent with previous receiver name receiver1 for GolintTest"
