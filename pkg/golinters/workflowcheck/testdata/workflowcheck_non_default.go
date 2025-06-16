//golangcitest:args -Eworkflowcheck
//golangcitest:config_path testdata/workflowcheck_non_default.yml
package testdata

import (
	"fmt"
	"time"
)

func CustomNonDeterministicFunc() { // want CustomNonDeterministicFunc:"declared non-deterministic"
	fmt.Println("This is marked as non-deterministic")
}

func CallsCustomNonDeterministicFunc() { // want CallsCustomNonDeterministicFunc:"calls non-deterministic function testdata.CustomNonDeterministicFunc"
	CustomNonDeterministicFunc()
}

func WrapperThatAcceptsNonDeterministic(fn func()) {
	fn()
}

func SafeCallToNonDeterministic() {
	WrapperThatAcceptsNonDeterministic(func() {
		time.Now()
	})
}
