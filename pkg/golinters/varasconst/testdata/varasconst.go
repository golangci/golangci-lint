//golangcitest:args -Evarasconst
package testdata

var Global_1 = ""

var (
	Global_2 = ""
	Global_3 = ""
	global_4 = ""
)

// const
var Global_Const_1 = ""

// const
var (
	Global_Const_2 = ""
	Global_Const_3 = ""
	global_Const_4 = ""
)

var (

	// const
	Global_Const_5 = ""

	Global_5 = ""

	// const
	Global_Const_6 = ""
)

func Direct_modification() {
	Global_1 = "modified"
	Global_2 = "modified"
	Global_3 = "modified"
	global_4 = "modified"
	Global_5 = "modified"

	Global_Const_1 = "modified" // want "assignment to global variable marked with const"
	Global_Const_2 = "modified" // want "assignment to global variable marked with const"
	Global_Const_3 = "modified" // want "assignment to global variable marked with const"
	global_Const_4 = "modified" // want "assignment to global variable marked with const"
	Global_Const_5 = "modified" // want "assignment to global variable marked with const"
	Global_Const_6 = "modified" // want "assignment to global variable marked with const"

	_ = global_4
	_ = global_Const_4
}

func Hide_global_var_by_locally_defined_one() {
	Global_1 := "defined"
	Global_1 = "modified"

	Global_Const_1 := "defined"
	Global_Const_1 = "modified"

	Global_2, Global_3 := "define multiple", "define multiple"
	Global_2 = "modified"
	Global_3 = "modified"

	Global_Const_2, Global_Const_3 := "define multiple", "define multiple"
	Global_Const_2 = "modified"
	Global_Const_3 = "modified"

	var global_4, global_Const_4 string
	global_4 = "modified"
	global_Const_4 = "modified"

	_ = Global_1
	_ = Global_2
	_ = Global_3
	_ = global_4
	_ = Global_Const_1
	_ = Global_Const_2
	_ = Global_Const_3
	_ = global_Const_4
}

func Assignment_in_if_stmt() {

	if Global_1 = "assignment in if condition"; true {
	}

	if Global_Const_1 = "assignment in if condition"; true { // want "assignment to global variable marked with const"
	}

	if true {
		Global_1 = "assignment in if body"
		Global_Const_1 = "assignment in if body" // want "assignment to global variable marked with const"
	}

	_ = Global_1
	_ = Global_Const_1
}

func Hidden_in_if_stmt_but_not_after() {
	if Global_1 := "assignment in if condition"; true {
		Global_1 = "modified"
		_ = Global_1
	}
	Global_1 = "modified"

	if Global_Const_1 := "assignment in if condition"; true {
		Global_Const_1 = "modified"
		_ = Global_Const_1
	}
	Global_Const_1 = "modified" // want "assignment to global variable marked with const"

	if true {
		Global_1 := "assignment in if body"
		Global_Const_1 := "assignment in if body"

		_ = Global_1
		_ = Global_Const_1
	}
	Global_1 = "modified"
	Global_Const_1 = "modified" // want "assignment to global variable marked with const"

	_ = Global_1
	_ = Global_Const_1
}

func Inside_func_literal() {

	func() {
		Global_1 = "modified"
		Global_Const_1 = "modified" // want "assignment to global variable marked with const"
	}()

	func1 := func() {
		Global_1 = "modified"
		Global_Const_1 = "modified" // want "assignment to global variable marked with const"
	}
	func1()
}
