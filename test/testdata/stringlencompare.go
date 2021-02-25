//args: -Estringlencompare
package testdata

func stringlencompare() {
	str := ""
	if len(str) == 0 { // ERROR "Compare string with \"\", don't compare len with 0"
	}
	if 0 != len(str) { // ERROR "Compare string with \"\", don't compare len with 0"
	}
	if len(returnString()+"foo") > 0 { // ERROR "Compare string with \"\", don't compare len with 0"
	}
}

func returnString() string { return "foo" }
