//golangcitest:args -Egoconst
//golangcitest:config_path testdata/goconst_find_duplicates.yml
package testdata

const AConst = "test"
const (
	AnotherConst   = "test" // want "const definition is duplicate of `AConst` at goconst_find_duplicates.go:5:7"
	UnrelatedConst = "i'm unrelated"
)

func Bazoo() {
	const Duplicated = "test" // want "const definition is duplicate of `AConst` at goconst_find_duplicates.go:5:7"

	const NotDuplicated = "i'm not duplicated"
}
