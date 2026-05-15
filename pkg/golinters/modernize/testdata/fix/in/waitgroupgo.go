//go:build go1.25

//golangcitest:args -Emodernize
//golangcitest:expected_exitcode 0
package waitgroup

import (
	"fmt"
	"sync"
)

// supported case for pattern 1.
func _() {
	var wg sync.WaitGroup
	wg.Add(1) // want "Goroutine creation can be simplified using WaitGroup.Go"
	go func() {
		defer wg.Done()
		fmt.Println()
	}()

	wg.Add(1) // want "Goroutine creation can be simplified using WaitGroup.Go"
	go func() {
		defer wg.Done()
	}()

	for range 10 {
		wg.Add(1) // want "Goroutine creation can be simplified using WaitGroup.Go"
		go func() {
			defer wg.Done()
			fmt.Println()
		}()
	}
}

// supported case for pattern 2.
func _() {
	var wg sync.WaitGroup
	wg.Add(1) // want "Goroutine creation can be simplified using WaitGroup.Go"
	go func() {
		fmt.Println()
		wg.Done()
	}()

	wg.Add(1) // want "Goroutine creation can be simplified using WaitGroup.Go"
	go func() {
		wg.Done()
	}()

	for range 10 {
		wg.Add(1) // want "Goroutine creation can be simplified using WaitGroup.Go"
		go func() {
			fmt.Println()
			wg.Done()
		}()
	}
}

// this function puts some wrong usages but waitgroup modernizer will still offer fixes.
func _() {
	var wg sync.WaitGroup
	wg.Add(1) // want "Goroutine creation can be simplified using WaitGroup.Go"
	go func() {
		defer wg.Done()
		defer wg.Done()
		fmt.Println()
	}()

	wg.Add(1) // want "Goroutine creation can be simplified using WaitGroup.Go"
	go func() {
		defer wg.Done()
		fmt.Println()
		wg.Done()
	}()

	wg.Add(1) // want "Goroutine creation can be simplified using WaitGroup.Go"
	go func() {
		fmt.Println()
		wg.Done()
		wg.Done()
	}()
}

// this function puts the unsupported cases of pattern 1.
func _() {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {}()

	wg.Add(1)
	go func(i int) {
		defer wg.Done()
		fmt.Println(i)
	}(1)

	wg.Add(1)
	go func() {
		fmt.Println()
		defer wg.Done()
	}()

	wg.Add(1)
	go func() { // noop: no wg.Done call inside function body.
		fmt.Println()
	}()

	go func() { // noop: no Add call before this go stmt.
		defer wg.Done()
		fmt.Println()
	}()

	wg.Add(2) // noop: only support Add(1).
	go func() {
		defer wg.Done()
	}()

	var wg1 sync.WaitGroup
	wg1.Add(1) // noop: Add and Done should be the same object.
	go func() {
		defer wg.Done()
		fmt.Println()
	}()

	wg.Add(1) // noop: Add and Done should be the same object.
	go func() {
		defer wg1.Done()
		fmt.Println()
	}()
}

// this function puts the unsupported cases of pattern 2.
func _() {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		wg.Done()
		fmt.Println()
	}()

	go func() { // noop: no Add call before this go stmt.
		fmt.Println()
		wg.Done()
	}()

	var wg1 sync.WaitGroup
	wg1.Add(1) // noop: Add and Done should be the same object.
	go func() {
		fmt.Println()
		wg.Done()
	}()

	wg.Add(1) // noop: Add and Done should be the same object.
	go func() {
		fmt.Println()
		wg1.Done()
	}()
}

type Server struct {
	wg sync.WaitGroup
}

type ServerContainer struct {
	serv Server
}

func _() {
	var s Server
	s.wg.Add(1) // want "Goroutine creation can be simplified using WaitGroup.Go"
	go func() {
		print()
		s.wg.Done()
	}()

	var sc ServerContainer
	sc.serv.wg.Add(1) // want "Goroutine creation can be simplified using WaitGroup.Go"
	go func() {
		print()
		sc.serv.wg.Done()
	}()

	var wg sync.WaitGroup
	arr := [1]*sync.WaitGroup{&wg}
	arr[0].Add(1) // want "Goroutine creation can be simplified using WaitGroup.Go"
	go func() {
		print()
		arr[0].Done()
	}()
}
