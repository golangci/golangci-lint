//args: -Eunused
package testdata

func fn1() {} // ERROR "func `fn1` is unused"

//nolint:unused
func fn2() { fn3() }

func fn3() {} // ERROR "func `fn3` is unused"

func fn4() { fn5() } // ERROR "func `fn4` is unused"

func fn5() {} // ERROR "func `fn5` is unused"

func fn6() { fn4() } // ERROR "func `fn6` is unused"

type unusedStruct struct{} // ERROR "type `unusedStruct` is unused"

type unusedStructNolintUnused struct{} //nolint:unused

type unusedStructNolintMegacheck struct{} //nolint:megacheck
