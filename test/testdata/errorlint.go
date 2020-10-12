//args: -Eerrorlint
package testdata

import (
	"errors"
	"log"
)

var errFoo = errors.New("foo")

func doThing() error {
	return errFoo
}

func compare() {
	err := doThing()
	if errors.Is(err, errFoo) {
		log.Println("ErrFoo")
	}
	if err == nil {
		log.Println("nil")
	}
	if err != nil {
		log.Println("nil")
	}
	if nil == err {
		log.Println("nil")
	}
	if nil != err {
		log.Println("nil")
	}
	if err == errFoo { // ERROR "comparing with == will fail on wrapped errors. Use errors.Is to check for a specific error"
		log.Println("errFoo")
	}
	if err != errFoo { // ERROR "comparing with != will fail on wrapped errors. Use errors.Is to check for a specific error"
		log.Println("not errFoo")
	}
	if errFoo == err { // ERROR "comparing with == will fail on wrapped errors. Use errors.Is to check for a specific error"
		log.Println("errFoo")
	}
	if errFoo != err { // ERROR "comparing with != will fail on wrapped errors. Use errors.Is to check for a specific error"
		log.Println("not errFoo")
	}
	switch err { // ERROR "switch on an error will fail on wrapped errors. Use errors.Is to check for specific errors"
	case errFoo:
		log.Println("errFoo")
	}
	switch doThing() { // ERROR "switch on an error will fail on wrapped errors. Use errors.Is to check for specific errors"
	case errFoo:
		log.Println("errFoo")
	}
}

type myError struct{}

func (*myError) Error() string {
	return "foo"
}

func doAnotherThing() error {
	return &myError{}
}

func typeCheck() {
	err := doAnotherThing()
	var me *myError
	if errors.As(err, &me) {
		log.Println("myError")
	}
	_, ok := err.(*myError) // ERROR "type assertion on error will fail on wrapped errors. Use errors.As to check for specific errors"
	if ok {
		log.Println("myError")
	}
	switch err.(type) { // ERROR "type switch on error will fail on wrapped errors. Use errors.As to check for specific errors"
	case *myError:
		log.Println("myError")
	}
	switch doAnotherThing().(type) { // ERROR "type switch on error will fail on wrapped errors. Use errors.As to check for specific errors"
	case *myError:
		log.Println("myError")
	}
	switch t := err.(type) { // ERROR "type switch on error will fail on wrapped errors. Use errors.As to check for specific errors"
	case *myError:
		log.Println("myError", t)
	}
	switch t := doAnotherThing().(type) { // ERROR "type switch on error will fail on wrapped errors. Use errors.As to check for specific errors"
	case *myError:
		log.Println("myError", t)
	}
}
