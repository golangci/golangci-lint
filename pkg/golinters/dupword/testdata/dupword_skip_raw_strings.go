//golangcitest:args -Edupword
//golangcitest:config_path testdata/dupword_skip_raw_strings.yml
package testdata

import "fmt"

func duplicateWordInRawString() {
	// this line include duplicated word the the // want `Duplicate words \(the\) found`
	s := `this line include duplicated word the the` // skip-raw-strings should skip this
	fmt.Println(s)
}
