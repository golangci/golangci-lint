package testdata

var nolintSpecific int    //nolint:gofmt
var nolintSpace int       // nolint: gofmt
var nolintSpaces int      //nolint:gofmt, govet
var nolintAll int         // nolint:all
var nolintAndAppendix int // nolint:all // another comment

//nolint:all
var nolintVarByPrecedingComment int

//nolint:all

var dontNolintVarByPrecedingCommentBecauseOfNewLine int

var nolintPrecedingVar string //nolint:all
var dontNolintVarByPrecedingCommentBecauseOfDifferentColumn int

//nolint:all
func nolintFuncByPrecedingComment() *string {
	xv := "v"
	return &xv
}

//nolint:all
// second line
func nolintFuncByPrecedingMultilineComment1() *string {
	xv := "v"
	return &xv
}

// first line
//nolint:all
func nolintFuncByPrecedingMultilineComment2() *string {
	xv := "v"
	return &xv
}

// first line
//nolint:all
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
