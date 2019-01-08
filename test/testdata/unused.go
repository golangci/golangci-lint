//args: -Eunused
package testdata

type unusedStruct struct{} // ERROR "type `unusedStruct` is unused"

type unusedStructNolintUnused struct{} //nolint:unused

type unusedStructNolintMegacheck struct{} //nolint:megacheck
