package testdata

var nolintSpecific int    // nolint:gofmt
var nolintSpace int       // nolint: gofmt
var nolintSpaces int      //nolint: gofmt, govet
var nolintAll int         // nolint
var nolintAndAppendix int // nolint // another comment

//nolint
var nolintVarByPrecedingComment int

//nolint

var dontNolintVarByPrecedingCommentBecauseOfNewLine int

var nolintPrecedingVar string //nolint
var dontNolintVarByPrecedingCommentBecauseOfDifferentColumn int

//nolint
func nolintFuncByPrecedingComment() *string {
	xv := "v"
	return &xv
}

//nolint
// second line
func nolintFuncByPrecedingMultilineComment1() *string {
	xv := "v"
	return &xv
}

// first line
//nolint
func nolintFuncByPrecedingMultilineComment2() *string {
	xv := "v"
	return &xv
}

// first line
//nolint
// third line
func nolintFuncByPrecedingMultilineComment3() *string {
	xv := "v"
	return &xv
}

var nolintAliasGAS bool //nolint:gas

var nolintAliasGosec bool //nolint:gosec

var nolintAliasUpperCase int // nolint: GAS

//nolint:errcheck
var (
	nolintVarBlockVar1 int
	nolintVarBlockVar2 int
)
