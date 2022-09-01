//golangcitest:args -Egofmt
//golangcitest:config_path testdata/configs/gofmt_rewrite_rules.yml
package testdata

import "fmt"

func GofmtRewriteRule() {
	vals := make([]int, 0)

	vals = append(vals, 1)
	vals = append(vals, 2)
	vals = append(vals, 3)

	slice := vals[1:len(vals)] // want "^File is not `gofmt`-ed"

	fmt.Println(slice)
}
