//golangcitest:args -Efuncreturn
package testdata

var a = 1

func funcreturnA() {

}
func funcreturnB() { // want "There is no newline before function"
}

// comment
var b = 2

// comment
func funcreturnC() {

}

func funcreturnD() {}
func funcreturnE() {}

func funcreturnF() {

}
func funcreturnG() { // want "There is no newline before function"

}
