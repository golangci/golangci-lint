//golangcitest:args -Echeckdeprecated
//golangcitest:config_path testdata/configs/checkdeprecated.yml
package testdata

// Deprecated: VarDeprecatedComment by GenDecl ValueSpec
var VarDeprecateComment = ""

// DEPRECATED: vars1/2/3 by GenDecl ValueSpec
var (
	// deprecated. vars1 by ValueSpec
	vars1Comment = "" // want "malformed deprecated header: the proper format is `Deprecated: <text>`"
	vars2Comment = "" // want "malformed deprecated header: use `Deprecated: ` \\(note the casing\\) instead of `DEPRECATED: `"
	vars3Comment = "" // want "malformed deprecated header: use `Deprecated: ` \\(note the casing\\) instead of `DEPRECATED: `"
)

// ConstDeprecated
// it's deprecated. ConstDeprecated by GenDecl ValueSpec
const ConstDeprecatedComment = "" // want "malformed deprecated header: the proper format is `Deprecated: <text>`"

// NOTE: deprecated. consts 1/2/3 by GenDecl ValueSpec
const (
	// deprecated, consts1 by ValueSpec
	consts1Comment = iota // want "malformed deprecated header: the proper format is `Deprecated: <text>`"
	consts2Comment        // want "malformed deprecated header: the proper format is `Deprecated: <text>`"
	consts3Comment        // want "malformed deprecated header: the proper format is `Deprecated: <text>`"
)

// FuncDeprecated
//
// Deprecated: don't use FuncDeprecated by FuncDecl
func FuncDeprecatedComment() {
}

type StructComment struct{}

// Deprecated, don't use it
func (p StructComment) StructFun() {} // want "malformed deprecated header: the proper format is `Deprecated: <text>`"

// Deprecated: use s3 instead of StructDeprecated, by GenDecl TypeSpec
type StructDeprecatedComment struct{}

func (p StructDeprecatedComment) Fun() {} // want "using deprecated: use s3 instead of StructDeprecated, by GenDecl TypeSpec"

// Deprecated.
type StructDeprecated2Comment struct{} // want "malformed deprecated header: the proper format is `Deprecated: <text>`"

// InterfaceDeprecated
//
// Deprecated, InterfaceDeprecated by GenDecl TypeSpec
type InterfaceDeprecatedComment interface{} // want "malformed deprecated header: the proper format is `Deprecated: <text>`"

// Deprecated struct 1/2/3 by GenDecl TypeSpec
type (
	// Deprecated struct1 by TypeSpec
	struct1Comment struct{} // want "malformed deprecated header: the proper format is `Deprecated: <text>`"
	struct2Comment struct { // want "malformed deprecated header: the proper format is `Deprecated: <text>`"
		F1 string
		// Deprecated F2 by Field
		F2 string // want "malformed deprecated header: the proper format is `Deprecated: <text>`"
	}
	struct3Comment struct{} // want "malformed deprecated header: the proper format is `Deprecated: <text>`"
)

func (s struct3Comment) fun1() {} // want "using deprecated: struct 1/2/3 by GenDecl TypeSpec"

// Deprecated fun2 by FuncDecl
func (s struct3Comment) fun2() {} // want "using deprecated: struct 1/2/3 by GenDecl TypeSpec"

// Deprecated interface 1/2/3
// by GenDecl TypeSpec
type (
	// Deprecated interface1 by TypeSpec
	interface1Comment interface{} // want "malformed deprecated header: the proper format is `Deprecated: <text>`"
	interface2Comment interface { // want "malformed deprecated header: the proper format is `Deprecated: <text>`"
	}
	interface3Comment interface { // want "malformed deprecated header: the proper format is `Deprecated: <text>`"
		fun1()
		// deprecated. interface3 fun2 by Field
		fun2() // want "malformed deprecated header: the proper format is `Deprecated: <text>`"
	}
)
