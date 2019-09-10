//args: -Egodox
//config: linters-settings.godox.keywords=FIXME,TODO
package testdata

func todoLeftInCode() {
	// TODO implement me // ERROR godox.go:6: Line contains FIXME/TODO: "TODO implement me"
	//TODO no space // ERROR godox.go:7: Line contains FIXME/TODO: "TODO no space"
	// TODO(author): 123 // ERROR godox.go:8: Line contains FIXME/TODO: "TODO\(author\): 123 // ERROR godox.go:8: L..."
	//TODO(author): 123 // ERROR godox.go:9: Line contains FIXME/TODO: "TODO\(author\): 123 // ERROR godox.go:9: L..."
	//TODO(author) 456 // ERROR godox.go:10: Line contains FIXME/TODO: "TODO\(author\) 456 // ERROR godox.go:10: L..."
	// TODO: qwerty // ERROR godox.go:11: Line contains FIXME/TODO: "TODO: qwerty // ERROR godox.go:11: Line ..."
	// todo 789 // ERROR godox.go:12: Line contains FIXME/TODO: "todo 789"
}
