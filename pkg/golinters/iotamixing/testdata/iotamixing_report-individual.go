//golangcitest:args -Eiotamixing
//golangcitest:config_path testdata/iotamixing_report-individual.yml
package testdata

import "fmt"

const (
	InvalidPerIndividualIotaDeclAboveAnything = "anything" // want "InvalidPerIndividualIotaDeclAboveAnything is a const with r-val in same const block as iota. keep iotas in separate const blocks"
	InvalidPerIndividualIotaDeclAboveNotZero  = iota
	InvalidPerIndividualIotaDeclAboveNotOne
	InvalidPerIndividualIotaDeclAboveNotTwo
)

const (
	InvalidPerIndividualIotaDeclBelowZero = iota
	InvalidPerIndividualIotaDeclBelowOne
	InvalidPerIndividualIotaDeclBelowTwo
	InvalidPerIndividualIotaDeclBelowAnything = "anything" // want "InvalidPerIndividualIotaDeclBelowAnything is a const with r-val in same const block as iota. keep iotas in separate const blocks"
)

const (
	InvalidPerIndividualIotaDeclBetweenZero = iota
	InvalidPerIndividualIotaDeclBetweenOne
	InvalidPerIndividualIotaDeclBetweenAnything = "anything" // want "InvalidPerIndividualIotaDeclBetweenAnything is a const with r-val in same const block as iota. keep iotas in separate const blocks"
	InvalidPerIndividualIotaDeclBetweenNotTwo
)

const (
	InvalidPerIndividualIotaDeclMultipleAbove   = "above" // want "InvalidPerIndividualIotaDeclMultipleAbove is a const with r-val in same const block as iota. keep iotas in separate const blocks"
	InvalidPerIndividualIotaDeclMultipleNotZero = iota
	InvalidPerIndividualIotaDeclMultipleNotOne
	InvalidPerIndividualIotaDeclMultipleBetween = "between" // want "InvalidPerIndividualIotaDeclMultipleBetween is a const with r-val in same const block as iota. keep iotas in separate const blocks"
	InvalidPerIndividualIotaDeclMultipleNotTwo
	InvalidPerIndividualIotaDeclMultipleBelow = "below" // want "InvalidPerIndividualIotaDeclMultipleBelow is a const with r-val in same const block as iota. keep iotas in separate const blocks"
)

const (
	ValidPerIndividualIotaZero = iota
	ValidPerIndividualIotaOne
	ValidPerIndividualIotaTwo
)

const (
	ValidPerIndividualRegularSomething = "something"
	ValidPerIndividualRegularAnything  = "anything"
)

func _() {
	fmt.Println("using the std import so goland doesn't nuke it")
}
