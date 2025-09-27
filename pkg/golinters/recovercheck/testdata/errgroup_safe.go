package testdata

import (
	"fmt"
	"log"
	"testdata/pkg"

	"golang.org/x/sync/errgroup"
)

// SafeErrgroup demonstrates safe errgroup usage with panic recovery
func SafeErrgroup() {
	var g errgroup.Group

	// This should NOT be flagged - safe errgroup goroutine with defer recover
	g.Go(func() error {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("Recovered from panic in errgroup: %v", r)
			}
		}()
		panic("This panic is recovered")
		return nil
	})

	// Another safe errgroup goroutine using a recovery function
	g.Go(func() error {
		defer recoverFromPanic()
		panic("Another recovered panic")
		return nil
	})

	g.Wait()
}

// SafeErrgroupWithExternalRecover demonstrates safe errgroup usage with external recovery function
func SafeErrgroupWithExternalRecover() {
	var g errgroup.Group

	// This should NOT be flagged - uses external recovery function
	g.Go(func() error {
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("External recovery: %v\n", r)
			}
		}()
		panic("This panic is recovered by external function")
		return nil
	})

	g.Wait()
}

// recoverFromPanic is a helper function that provides panic recovery
func recoverFromPanic() {
	if r := recover(); r != nil {
		log.Printf("Recovered from panic: %v", r)
	}
}

// SafeErrgroupWithNamedRecover demonstrates safe errgroup usage with named recovery function
func SafeErrgroupWithNamedRecover() {
	var g errgroup.Group

	// This should NOT be flagged - uses named recovery function
	g.Go(func() error {
		defer recoverFromPanic()
		panic("This panic is recovered by named function")
		return nil
	})

	g.Wait()
}

// SafeErrgroupWithRecoverFromPackage demonstrates safe errgroup usage with recovery from another package
func SafeErrgroupWithRecoverFromPackage() {
	var g errgroup.Group

	// This should NOT be flagged - uses recovery from another package
	g.Go(func() error {
		defer pkg.PanicRecover()
		panic("This panic is recovered by another package")
		return nil
	})

	g.Wait()
}
