// Copyright (c) 2025 golangci-lint <someone@example.com>. // want "invalid license header"
// This file is a part of golangci-lint.

//golangcitest:args -Egolicenser
//golangcitest:config_path testdata/golicenser.yml
package testdata

import "fmt"

func main() {
	fmt.Println("Hello, world!")
}
