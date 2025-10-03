//golangcitest:args -Erecovercheck
package testdata

import (
	"log"
	"testdata/pkg"
	aliaspkg "testdata/pkg"
	"time"
)

// UnsafeFunction demonstrates unsafe goroutines without panic recovery
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

// SafeGoroutine2 uses a recovery function defined in the same file
func SafeGoroutine2() {
	go func() {
		defer sameFileRecover()
		panic("oh no")
	}()
}

// SafeGoroutine3 uses a recovery function with any name
func SafeGoroutine3() {
	go func() {
		defer anyName()
		panic("oh no")
	}()
}

// Group - Mock errgroup for testing
type Group struct{}

func (g *Group) Go(f func() error) {
	go f() // want "goroutine created without panic recovery"
}

func (g *Group) Wait() error {
	return nil
}

// UnsafeErrgroup demonstrates unsafe errgroup usage without panic recovery
func UnsafeErrgroup() {
	var g Group

	// This should be flagged - unsafe errgroup goroutine
	g.Go(func() error { // want "errgroup goroutine created without panic recovery"
		panic("This will crash the program")
	})

	// Another unsafe errgroup goroutine
	g.Go(func() error { // want "errgroup goroutine created without panic recovery"
		panic("Another unrecovered panic")
	})

	g.Wait()
}

// SafeErrgroup demonstrates safe errgroup usage with panic recovery
func SafeErrgroup() {
	var g Group

	// This should NOT be flagged - safe errgroup goroutine with defer recover
	g.Go(func() error {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("Recovered from panic in errgroup: %v", r)
			}
		}()
		panic("This panic is recovered")
	})

	// Another safe errgroup goroutine using a recovery function
	g.Go(func() error {
		defer recoverFromPanic()
		panic("Another recovered panic")
	})

	g.Wait()
}

// recoverFromPanic is a helper function that provides panic recovery
func recoverFromPanic() {
	if r := recover(); r != nil {
		log.Printf("Recovered from panic: %v", r)
	}
}

// SafeGoroutineWithExternalRecover uses a recovery function from an external package
func SafeGoroutineWithExternalRecover() {
	go func() {
		defer pkg.PanicRecover()
		panic("oh no")
	}()
}

// SafeGoroutineWithAliasImport uses a recovery function from another package with import alias
func SafeGoroutineWithAliasImport() {
	go func() {
		defer aliaspkg.PanicRecover()
		panic("oh no")
	}()
}

// SafeErrgroupWithExternalRecover demonstrates safe errgroup usage with external recovery function
func SafeErrgroupWithExternalRecover() {
	var g Group

	// This should NOT be flagged - uses external recovery function
	g.Go(func() error {
		defer pkg.PanicRecover()
		panic("This panic is recovered by external package")
	})

	g.Wait()
}
