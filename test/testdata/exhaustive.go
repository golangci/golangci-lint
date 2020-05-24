//args: -Eexhaustive
package testdata

import (
	"github.com/golangci/golangci-lint/test/testdata/exhaustive"
)

func processDirection(d exhaustive.Direction) {
	switch d { // ERROR "missing cases in switch of type exhaustive.Direction: East, West"
	case North, South:
	}
}
