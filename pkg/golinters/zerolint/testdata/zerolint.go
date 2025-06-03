//golangcitest:args -Ezerolint
//golangcitest:config_path testdata/zerolint.yml
package testdata

import (
	"errors"
	"fmt"
)

type MyError struct{}

func (*MyError) Error() string { // want "zl:err"
	return "my error"
}

func DoWork() error {
	return &MyError{}
}

type Excluded struct{}

func main() {
	if err := DoWork(); errors.Is(err, &MyError{}) { // want "zl:cme"
		fmt.Println(`Got "my error"`)
	}

	_ = &Excluded{} == &Excluded{}
}
