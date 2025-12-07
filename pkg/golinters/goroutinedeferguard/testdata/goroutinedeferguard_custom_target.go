//golangcitest:args -Egoroutinedeferguard
//golangcitest:config_path testdata/goroutinedeferguard_custom_target.yml
package testdata

import (
	"fmt"

	wrong "github.com/golangci/golangci-lint/v2/pkg/golinters/goroutinedeferguard/testdata/other"
	"github.com/golangci/golangci-lint/v2/pkg/golinters/goroutinedeferguard/testdata/right"
)

func goodAnonymous() {
	go func() {
		defer right.MyPanicHandler()
		fmt.Println("Hello, World!")
	}()
}

func badAnonymous() {
	go func() { // want "missing defer call to custompattern/right.MyPanicHandler"
		fmt.Println("Hello, World!")
	}()
}

func badWrongPackage() {
	go func() { // want "missing defer call to custompattern/right.MyPanicHandler"
		defer wrong.MyPanicHandler()
		fmt.Println("Hello, World!")
	}()
}
