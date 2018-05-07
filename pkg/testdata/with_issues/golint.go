package testdata

var Go_lint string // ERROR "don't use underscores in Go names; var Go_lint should be GoLint"

func ExportedFuncWithNoComment() {
}

var ExportedVarWithNoComment string

type ExportedStructWithNoComment struct{}

type ExportedInterfaceWithNoComment interface{}

// Bad comment // ERROR "comment on exported function ExportedFuncWithBadComment should be of the form .ExportedFuncWithBadComment \.\.\.."
func ExportedFuncWithBadComment() {}
