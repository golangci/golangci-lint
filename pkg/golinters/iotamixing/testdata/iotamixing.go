//golangcitest:args -Eiotamixing
package testdata

import "fmt"

// iota mixing in const block containing an iota and r-val declared above.
const ( // want "iota mixing. keep iotas in separate blocks to consts with r-val"
	InvalidPerBlockIotaDeclAboveAnything = "anything"
	InvalidPerBlockIotaDeclAboveNotZero  = iota
	InvalidPerBlockIotaDeclAboveNotOne
	InvalidPerBlockIotaDeclAboveNotTwo
)

// iota mixing in const block containing an iota and r-val declared below.
const ( // want "iota mixing. keep iotas in separate blocks to consts with r-val"
	InvalidPerBlockIotaDeclBelowZero = iota
	InvalidPerBlockIotaDeclBelowOne
	InvalidPerBlockIotaDeclBelowTwo
	InvalidPerBlockIotaDeclBelowAnything = "anything"
)

// iota mixing in const block containing an iota and r-val declared between consts.
const ( // want "iota mixing. keep iotas in separate blocks to consts with r-val"
	InvalidPerBlockIotaDeclBetweenZero = iota
	InvalidPerBlockIotaDeclBetweenOne
	InvalidPerBlockIotaDeclBetweenAnything = "anything"
	InvalidPerBlockIotaDeclBetweenNotTwo
)

// iota mixing in const block containing an iota and r-vals declared above, between, and below consts.
const ( // want "iota mixing. keep iotas in separate blocks to consts with r-val"
	InvalidPerBlockIotaDeclMultipleAbove   = "above"
	InvalidPerBlockIotaDeclMultipleNotZero = iota
	InvalidPerBlockIotaDeclMultipleNotOne
	InvalidPerBlockIotaDeclMultipleBetween = "between"
	InvalidPerBlockIotaDeclMultipleNotTwo
	InvalidPerBlockIotaDeclMultipleBelow = "below"
)

// no iota mixing in a const block containing an iota and no r-vals.
const (
	ValidPerBlockIotaZero = iota
	ValidPerBlockIotaOne
	ValidPerBlockIotaTwo
)

// no iota mixing in a const block containing r-vals and no iota.
const (
	ValidPerBlockRegularSomething = "something"
	ValidPerBlockRegularAnything  = "anything"
)

func _() {
	fmt.Println("using the std import so goland doesn't nuke it")
}
