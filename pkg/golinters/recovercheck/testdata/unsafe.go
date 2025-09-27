package testdata

import "time"

func UnsafeFunction() {
	// This should be flagged - unsafe goroutine
	go func() { // want "goroutine created without panic recovery"
		panic("This will crash the program")
	}()

	// Another unsafe goroutine
	go func() { // want "goroutine created without panic recovery"
		time.Sleep(1 * time.Second)
		panic("Another unrecovered panic")
	}()
}
