//golangcitest:args -Emodernize
package testdata

func _(x interface{}) {} // want "interface{} can be replaced by any"

func _() {
	var x interface{} // want "interface{} can be replaced by any"
	const any = 1
	var y interface{} // nope: any is shadowed here
	_, _ = x, y
}
