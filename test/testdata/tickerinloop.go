//args: -Etickerinloop
package testdata

import "time"

func tickerInLoop() {
	for {
		t := time.NewTicker(1 * time.Second) // ERROR "ticker found in loop"
		t.Stop()
	}
}
