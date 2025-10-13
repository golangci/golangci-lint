//golangcitest:args -Emodernize
//golangcitest:expected_exitcode 0
package fmtappendf

import (
	"fmt"
)

func two() string {
	return "two"
}

func bye() {
	_ = fmt.Appendf(nil, "bye %d", 1) // want "Replace .*Sprintf.* with fmt.Appendf"
}

func funcsandvars() {
	one := "one"
	_ = fmt.Appendf(nil, "bye %d %s %s", 1, two(), one) // want "Replace .*Sprintf.* with fmt.Appendf"
}

func typealias() {
	type b = byte
	type bt = []byte
	_ = fmt.Appendf(nil, "bye %d", 1) // want "Replace .*Sprintf.* with fmt.Appendf"
	_ = fmt.Appendf(nil, "bye %d", 1) // want "Replace .*Sprintf.* with fmt.Appendf"
}

func otherprints() {
	_ = fmt.Append(nil, "bye %d", 1)   // want "Replace .*Sprint.* with fmt.Append"
	_ = fmt.Appendln(nil, "bye %d", 1) // want "Replace .*Sprintln.* with fmt.Appendln"
}

func comma() {
	type S struct{ Bytes []byte }
	var _ = struct{ A S }{
		A: S{
			Bytes: // want "Replace .*Sprint.* with fmt.Appendf"
			fmt.Appendf(nil, "%d", 0),
		},
	}
	_ = // want "Replace .*Sprint.* with fmt.Appendf"
		fmt.Appendf(nil, "%d", 0)
}
