package testdata

import (
	"log"

	"recovercheck/pkg"
)

// SafeGoroutine demonstrates a goroutine with proper panic recovery
func SafeGoroutine() {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Println("Recovered from panic:", r)
			}
		}()
		panic("oh no")
	}()
}

// sameFileRecover is a recovery function defined in the same file
func sameFileRecover() func() {
	return func() {
		if r := recover(); r != nil {
			log.Println("Recovered from panic:", r)
		}
	}
}

// anyName is a recovery function with any name -- to check that the analyzer doesn't rely on the name
func anyName() func() {
	return func() {
		if r := recover(); r != nil {
			log.Println("Recovered from panic:", r)
		}
	}
}

// SafeGoroutine2 uses a recovery function from another package
func SafeGoroutine2() {
	go func() {
		defer pkg.PanicRecover()
		panic("oh no")
	}()
}

// SafeGoroutine3 uses a recovery function defined in the same file
func SafeGoroutine3() {
	go func() {
		defer sameFileRecover()
		panic("oh no")
	}()
}

// SafeGoroutine4 uses a recovery function with any name
func SafeGoroutine4() {
	go func() {
		defer anyName()
		panic("oh no")
	}()
}
