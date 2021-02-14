//args: -Edurationcheck
package testdata

import "time"

func waitFor(someDuration time.Duration) {
	timeToWait := someDuration * time.Second // ERROR "Multiplication of durations: `someDuration \\* time.Second` "
	time.Sleep(timeToWait)
}
