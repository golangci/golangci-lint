package testdata

import (
	"golang.org/x/sync/errgroup"
)

// UnsafeErrgroup demonstrates unsafe errgroup usage without panic recovery
func UnsafeErrgroup() {
	var g errgroup.Group

	// This should be flagged - unsafe errgroup goroutine
	g.Go(func() error { // want "errgroup goroutine created without panic recovery"
		panic("This will crash the program")
		return nil
	})

	// Another unsafe errgroup goroutine
	g.Go(func() error { // want "errgroup goroutine created without panic recovery"
		panic("Another unrecovered panic")
		return nil
	})

	g.Wait()
}

// UnsafeErrgroupStdlib demonstrates unsafe errgroup usage with regular goroutines
func UnsafeErrgroupStdlib() {
	// This shows a pattern where someone manually creates goroutines instead of using errgroup
	go func() { // want "goroutine created without panic recovery"
		panic("This will crash the program")
	}()
}

// UnsafeErrgroupWithAlias demonstrates unsafe errgroup usage with import alias
func UnsafeErrgroupWithAlias() {
	// Using a type alias to simulate different import patterns
	type Group = errgroup.Group
	var g Group

	g.Go(func() error { // want "errgroup goroutine created without panic recovery"
		panic("This will crash the program")
		return nil
	})

	g.Wait()
}
