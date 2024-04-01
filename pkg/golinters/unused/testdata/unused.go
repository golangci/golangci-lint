//golangcitest:args -Eunused
package testdata

func fn1() {} // want "func `fn1` is unused"

//nolint:unused
func fn2() { fn3() }

func fn3() {} // want "func `fn3` is unused"

func fn4() { fn5() } // want "func `fn4` is unused"

func fn5() {} // want "func `fn5` is unused"

func fn6() { fn4() } // want "func `fn6` is unused"

type unusedStruct struct{} // want "type `unusedStruct` is unused"

type unusedStructNolintUnused struct{} //nolint:unused

type unusedStructNolintMegacheck struct{} //nolint:megacheck
