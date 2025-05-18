//golangcitest:args -Eembeddedstructfieldcheck
package testdata

import "time"

type ValidStructWithSingleLineComments struct {
	// time.Time Single line comment
	time.Time

	// version Single line comment
	version int
}

type StructWithSingleLineComments struct {
	// time.Time Single line comment
	time.Time // want `there must be an empty line separating embedded fields from regular fields`
	// version Single line comment
	version int
}

type StructWithMultiLineComments struct {
	// time.Time Single line comment
	time.Time // want `there must be an empty line separating embedded fields from regular fields`
	// version Single line comment
	// very long comment
	version int
}
