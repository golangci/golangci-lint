//golangcitest:args -Eembeddedstructfieldcheck
//golangcitest:expected_exitcode 0
package testdata

import (
	"time"
)

type ValidStruct struct {
	time.Time

	version int
}

type NoSpaceStruct struct {
	time.Time // want `there must be an empty line separating embedded fields from regular fields`
	version   int
}

type EmbeddedWithPointers struct {
	*time.Time // want `there must be an empty line separating embedded fields from regular fields`
	version    int
}
