package testdata

import (
	"fmt" // ERROR "File is not goimports-ed"
	"github.com/golangci/golangci-lint/pkg/config"
)

func bar() {
	fmt.Print("x")
	_ = config.Config{}
}
