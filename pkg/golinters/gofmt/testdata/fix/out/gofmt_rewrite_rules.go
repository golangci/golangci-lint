//golangcitest:args -Egofmt
//golangcitest:config_path testdata/gofmt_rewrite_rules.yml
//golangcitest:expected_exitcode 0
package p

import "fmt"

func GofmtRewriteRule() {
	values := make([]int, 0)

	values = append(values, 1)
	values = append(values, 2)
	values = append(values, 3)

	slice := values[1:]

	fmt.Println(slice)
}

func GofmtRewriteRule2() {
	var to any

	fmt.Println(to)
}
