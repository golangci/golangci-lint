// Copyright (c) 2025 golangci-lint <someone@example.com>.
// This file is a part of golangci-lint.

//golangcitest:args -Egolicenser
//golangcitest:expected_exitcode 0
//golangcitest:config_path testdata/golicenser-fix.yml
package testdata

import "fmt"

func main() {
	fmt.Println("Hello, world!")
}
