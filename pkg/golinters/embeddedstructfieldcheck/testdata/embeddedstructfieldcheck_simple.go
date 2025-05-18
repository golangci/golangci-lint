//golangcitest:args -Eembeddedstructfieldcheck
package testdata

import (
	"context"
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

type NotSortedStruct struct {
	version int

	time.Time // want `embedded fields should be listed before regular fields`
}

type MixedEmbeddedAndNotEmbedded struct {
	context.Context

	name string

	time.Time // want `embedded fields should be listed before regular fields`

	age int
}

type EmbeddedWithPointers struct {
	*time.Time // want `there must be an empty line separating embedded fields from regular fields`
	version    int
}
