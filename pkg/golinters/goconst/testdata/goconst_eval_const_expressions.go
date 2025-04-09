//golangcitest:args -Egoconst
//golangcitest:config_path testdata/goconst_eval_const_expressions.yml
package testdata

const (
	prefix = "example.com/"
	API    = prefix + "api"
	Web    = prefix + "web"
)

const Full = "example.com/api"

func _() {
	a0 := "example.com/api" // want "string `example.com/api` has 3 occurrences, but such constant `API` already exists"
	a1 := "example.com/api"
	a2 := "example.com/api"

	_ = a0
	_ = a1
	_ = a2

	b0 := "example.com/web" // want "string `example.com/web` has 3 occurrences, but such constant `Web` already exists"
	b1 := "example.com/web"
	b2 := "example.com/web"

	_ = b0
	_ = b1
	_ = b2
}
