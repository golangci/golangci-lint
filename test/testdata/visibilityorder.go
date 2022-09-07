//golangcitest:args -Evisibilityorder
package testdata

func visibilityorderTestunexported() { // ERROR
}

func visibilityorderTestExported() {
}
