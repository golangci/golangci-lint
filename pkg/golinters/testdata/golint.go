package testdata

var go_lint string // ERROR "don't use underscores in Go names; var go_lint should be goLint"

func ExportedFuncWithNoComment() {
}

var ExportedVarWithNoComment string

type ExportedStructWithNoComment struct{}

type ExportedInterfaceWithNoComment interface{}

// Bad comment // ERROR "comment on exported function ExportedFuncWithBadComment should be of the form .ExportedFuncWithBadComment \.\.\.."
func ExportedFuncWithBadComment() {}
