//golangcitest:args -Egoconst
//golangcitest:config_path testdata/goconst_eval_expressions.yml
package testdata

const Host = "www.golangci.com"
const LinterPath = Host + "/goconst"

const path = "www.golangci.com/goconst" // want "const definition is duplicate of `LinterPath` at goconst_eval_expressions.go:6:20"

const KiB = 1 << 10

func EvalExpr() {
	println(path)

	const duplicated = "www.golangci.com/goconst" // want "const definition is duplicate of `LinterPath` at goconst_eval_expressions.go:6:20"
	println(duplicated)

	var unique = "www.golangci.com/another-linter"
	println(unique)

	const Kilobytes = 1024 // want "const definition is duplicate of `KiB` at goconst_eval_expressions.go:10:13"
	println(Kilobytes)
}
