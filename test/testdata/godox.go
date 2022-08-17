//golangcitest:args -Egodox
//golangcitest:config linters-settings.godox.keywords=FIXME,TODO
package testdata

func todoLeftInCode() {
	// TODO implement me // ERROR `Line contains FIXME/TODO: "TODO implement me`
	//TODO no space // ERROR `Line contains FIXME/TODO: "TODO no space`
	// TODO(author): 123 // ERROR `Line contains FIXME/TODO: "TODO\(author\): 123`
	//TODO(author): 123 // ERROR `Line contains FIXME/TODO: "TODO\(author\): 123`
	//TODO(author) 456 // ERROR `Line contains FIXME/TODO: "TODO\(author\) 456`
	// TODO: qwerty // ERROR `Line contains FIXME/TODO: "TODO: qwerty`
	// todo 789 // ERROR `Line contains FIXME/TODO: "todo 789`
}
