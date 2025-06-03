//golangcitest:args -Ezerolint
//golangcitest:config_path testdata/zerolint-fix.yml
package main

import (
	"errors"
	"fmt"
)

type MyError struct{}

func (*MyError) Error() string {
	return "my error"
}

func DoWork() error {
	return &MyError{}
}

func main() {
	if err := DoWork(); errors.Is(err, &MyError{}) {
		fmt.Println(`Got "my error"`)
	}
}
