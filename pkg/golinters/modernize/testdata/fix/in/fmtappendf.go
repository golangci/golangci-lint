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
	_ = []byte(fmt.Sprintf("bye %d", 1)) // want "Replace .*Sprintf.* with fmt.Appendf"
}

func funcsandvars() {
	one := "one"
	_ = []byte(fmt.Sprintf("bye %d %s %s", 1, two(), one)) // want "Replace .*Sprintf.* with fmt.Appendf"
}

func typealias() {
	type b = byte
	type bt = []byte
	_ = []b(fmt.Sprintf("bye %d", 1)) // want "Replace .*Sprintf.* with fmt.Appendf"
	_ = bt(fmt.Sprintf("bye %d", 1))  // want "Replace .*Sprintf.* with fmt.Appendf"
}

func otherprints() {
	_ = []byte(fmt.Sprint("bye %d", 1))   // want "Replace .*Sprint.* with fmt.Append"
	_ = []byte(fmt.Sprintln("bye %d", 1)) // want "Replace .*Sprintln.* with fmt.Appendln"
}

func comma() {
	type S struct{ Bytes []byte }
	var _ = struct{ A S }{
		A: S{
			Bytes: []byte( // want "Replace .*Sprint.* with fmt.Appendf"
				fmt.Sprintf("%d", 0),
			),
		},
	}
	_ = []byte( // want "Replace .*Sprint.* with fmt.Appendf"
		fmt.Sprintf("%d", 0),
	)
}
