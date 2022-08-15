//golangcitest:args -Ecyclop
//golangcitest:config_path testdata/configs/cyclop.yml
package testdata

func cyclopComplexFunc(s string) { // ERROR "calculated cyclomatic complexity for function cyclopComplexFunc is 22, max is 15"
	if s == "1" || s == "2" || s == "3" || s == "4" || s == "5" || s == "6" || s == "7" {
		return
	}
	if s == "1" || s == "2" || s == "3" || s == "4" || s == "5" || s == "6" || s == "7" {
		return
	}
	if s == "1" || s == "2" || s == "3" || s == "4" || s == "5" || s == "6" || s == "7" {
		return
	}
}
