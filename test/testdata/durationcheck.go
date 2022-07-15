//golangcitest:args -Edurationcheck
package testdata

import (
	"fmt"
	"time"
)

type durationCheckData struct {
	i int
	d time.Duration
}

func durationcheckCase01() {
	dcd := durationCheckData{i: 10}
	_ = time.Duration(dcd.i) * time.Second
}

func durationcheckCase02() {
	dcd := durationCheckData{d: 10 * time.Second}
	_ = dcd.d * time.Second // ERROR "Multiplication of durations: `dcd.d \\* time.Second`"
}

func durationcheckCase03() {
	seconds := 10
	fmt.Print(time.Duration(seconds) * time.Second)
}

func durationcheckCase04(someDuration time.Duration) {
	timeToWait := someDuration * time.Second // ERROR "Multiplication of durations: `someDuration \\* time.Second`"
	time.Sleep(timeToWait)
}

func durationcheckCase05() {
	someDuration := 2 * time.Second
	timeToWait := someDuration * time.Second // ERROR "Multiplication of durations: `someDuration \\* time.Second`"
	time.Sleep(timeToWait)
}
