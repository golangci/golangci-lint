//golangcitest:args -Egoconcurrencylint
package testdata

import "sync"

var packageMu1 sync.Mutex
var packageWG1 sync.WaitGroup

func badPackageMutex() {
	packageMu1.Lock() // want "mutex 'packageMu1' is locked but not unlocked"
}

func badPackageWaitGroup() {
	packageWG1.Add(1) // want "waitgroup 'packageWG1' has Add without corresponding Done"
	packageWG1.Wait()
}

func badWaitGroupGoAfterWait() {
	var wg sync.WaitGroup
	wg.Wait()
	wg.Go(func() {}) // want "waitgroup 'wg' Go called after Wait"
}

func goodWaitGroupGo() {
	var wg sync.WaitGroup
	wg.Go(func() {})
	wg.Wait()
}
