//args: -Ecopywaitgroup
package testdata

import (
	"sync"
)

var globalWait sync.WaitGroup

func g(wg sync.WaitGroup) {
	wg.Done()
}

func pg(wg *sync.WaitGroup) {
	wg.Done()
}

func closure() {
	var wg sync.WaitGroup
	a := 1

	wg.Add(1)
	go func(a int, wg sync.WaitGroup) { // ERROR "`wg` of arg `2` is passed as a value. you must pass as a pointer. the goroutine will be deadlock!"
		println(a)
		wg.Done()
	}(a, wg)
	wg.Wait()

	wg.Add(1)
	go func(a int, wg *sync.WaitGroup) { // OK
		println(a)
		wg.Done()
	}(a, &wg)
	wg.Wait()

	wg.Add(1)
	go func() { // OK
		wg.Done()
	}()
	wg.Wait()
}

func declared() {
	var wg sync.WaitGroup
	wg.Add(2)
	go g(wg)   // ERROR "`wg` of arg `1` is passed as a value. you must pass as a pointer. the goroutine will be deadlock!"
	go pg(&wg) // OK
	wg.Wait()
}

func global() {
	go g(globalWait) // OK
}
