//golangcitest:args -Egoconst
//golangcitest:config_path testdata/goconst_find_duplicates.yml
package testdata

const SingleConst = "single constant"

const (
	GroupedConst1 = "grouped constant"
	GroupedConst2 = "another grouped"
)

const (
	GroupedDuplicateConst1 = "grouped duplicate value"
	GroupedDuplicateConst2 = "grouped duplicate value" // want "This constant is a duplicate of `GroupedDuplicateConst1` at .*goconst_find_duplicates.go:13:2"
)

const DuplicateConst1 = "duplicate value"

const DuplicateConst2 = "duplicate value" // want "This constant is a duplicate of `DuplicateConst1` at .*goconst_find_duplicates.go:17:7"

const (
	SpecialDuplicateConst1 = "special\nvalue\twith\rchars"
	SpecialDuplicateConst2 = "special\nvalue\twith\rchars" // want "This constant is a duplicate of `SpecialDuplicateConst1` at .*goconst_find_duplicates.go:22:2"
)

func _() {
	const DuplicateScopedConst1 = "duplicate scoped value"
	const DuplicateScopedConst2 = "duplicate scoped value" // want "This constant is a duplicate of `DuplicateScopedConst1` at .*goconst_find_duplicates.go:27:8"
}
