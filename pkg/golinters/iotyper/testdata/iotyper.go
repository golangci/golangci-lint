//golangcitest:args -Eiotyper
package testdata

// Basic iota without type
const (
	BasicWithoutType = iota // want "iota used without type specification"
)

// Single-line const declarations
const SingleLineWithoutType = iota  // want "iota used without type specification"
const SingleLineWithType int = iota // OK: has type

// Multiple constants with inherited iota
const (
	FirstInGroup  = iota // want "iota used without type specification"
	SecondInGroup        // No warning: inherits iota value but doesn't use iota directly
	ThirdInGroup         // No warning: same as above
)

// Non-iota constants (no warnings expected)
const (
	PlainNumber = 42
	PlainString = "hello"
	PlainBool   = true
)

// iota in expressions
const (
	IotaPlusOne  = iota + 1  // want "iota used without type specification"
	IotaShifted  = 1 << iota // want "iota used without type specification"
	IotaMultiple = iota * 2  // want "iota used without type specification"
)

// iota with explicit type
const (
	WithTypeInt      int = iota // OK: explicit int type
	WithTypeIntAgain int = iota // OK: explicit int type
)

// Mixed type specifications
const (
	MixedWithType      int = iota // OK: has type
	MixedWithoutType       = iota // want "iota used without type specification"
	MixedAgainWithType int = iota // OK: has type
)
