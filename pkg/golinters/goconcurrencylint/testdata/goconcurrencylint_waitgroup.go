//golangcitest:args -Egoconcurrencylint
package testdata

import "sync"

type workerPool struct {
	wg sync.WaitGroup
}

func handleExternalWorkForTest(wg *sync.WaitGroup) {
	defer wg.Done()
}

func runWithDoneCallbackForTest(done func()) {
	defer done()
}

func badExtraDoneCase() {
	var wg sync.WaitGroup
	wg.Add(1)
	wg.Done()
	wg.Done() // want "waitgroup 'wg' has Done without corresponding Add"
	wg.Wait()
}

func badAddAfterWaitCase() {
	var wg sync.WaitGroup
	wg.Wait()
	go func() {
		wg.Add(1) // want "waitgroup 'wg' Add called after Wait"
		wg.Done()
	}()
}

func badAddAfterWaitMainFlowCase() {
	var wg sync.WaitGroup
	wg.Wait()
	wg.Add(1) // want "waitgroup 'wg' Add called after Wait"
	wg.Done()
}

func badLoopMissingDoneCase() {
	var wg sync.WaitGroup
	for i := 0; i < 2; i++ {
		wg.Add(1) // want "waitgroup 'wg' has Add without corresponding Done"
		if i == 0 {
			go func() {
				wg.Done()
			}()
		}
	}
	wg.Wait()
}

func badPrematureReturnCase() {
	var wg sync.WaitGroup
	wg.Add(1) // want "waitgroup 'wg' has Add without corresponding Done"
	go func() {
		return
		wg.Done()
	}()
	wg.Wait()
}

func badMethodWaitGroupCase() {
	var wp workerPool
	wp.wg.Add(1) // want "waitgroup 'wp.wg' has Add without corresponding Done"
	wp.wg.Wait()
}

func badSwitchDefaultOnlyDoneCase() {
	var wg sync.WaitGroup
	wg.Add(1) // want "waitgroup 'wg' has Add without corresponding Done"
	go func() {
		x := 1
		switch x {
		case 2:
			// missing Done in this branch
		default:
			wg.Done()
		}
	}()
	wg.Wait()
}

func goodDeferredDoneCase() {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
	}()
	wg.Wait()
}

func goodSwitchWithDefaultCase() {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		x := 1
		switch x {
		case 2:
			wg.Done()
		default:
			wg.Done()
		}
	}()
	wg.Wait()
}

func goodWaitGroupPassedToHelperCase() {
	var wg sync.WaitGroup
	wg.Add(1)
	go handleExternalWorkForTest(&wg)
	wg.Wait()
}

func goodWaitGroupMethodPassedCase() {
	var wg sync.WaitGroup
	wg.Add(1)
	go runWithDoneCallbackForTest(wg.Done)
	wg.Wait()
}

func goodReuseWaitGroupCase() {
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
	}()
	wg.Wait()

	wg.Add(1)
	go func() {
		defer wg.Done()
	}()
	wg.Wait()
}
