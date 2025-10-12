//golangcitest:args -Emodernize
//golangcitest:config_path testdata/modernize_custom.yml
//golangcitest:expected_exitcode 0
package testdata

func _(x interface{}) {}

func _() {
	var x interface{}
	const any = 1
	var y interface{} // nope: any is shadowed here
	_, _ = x, y
}
