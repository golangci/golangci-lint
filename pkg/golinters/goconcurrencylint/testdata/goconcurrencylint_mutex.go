//golangcitest:args -Egoconcurrencylint
package testdata

import "sync"

type safeMap struct {
	mu sync.Mutex
	rw sync.RWMutex
}

func badDoubleUnlockCase() {
	var mu sync.Mutex
	mu.Lock()
	mu.Unlock()
	mu.Unlock() // want "mutex 'mu' is unlocked but not locked"
}

func badConditionalMissingUnlockCase() {
	var mu sync.Mutex
	if true {
		mu.Lock() // want "mutex 'mu' is locked but not unlocked in if"
	}
}

func badElseIfMissingUnlockCase() {
	var mu sync.Mutex
	cond1 := false
	cond2 := true

	if cond1 {
		mu.Lock()
		mu.Unlock()
	} else if cond2 {
		mu.Lock() // want "mutex 'mu' is locked but not unlocked in if"
	}
}

func badGoroutineMutexCase() {
	var mu sync.Mutex
	ch := make(chan struct{})
	go func() {
		mu.Lock() // want "mutex 'mu' is locked but not unlocked in goroutine"
		<-ch
	}()
}

func badStructFieldMutexCase() {
	var sm safeMap
	sm.mu.Lock() // want "mutex 'sm.mu' is locked but not unlocked"
}

func badGoroutineDeferUnlockWithoutLockCase() {
	var mu sync.Mutex
	ch := make(chan struct{})
	go func() {
		defer mu.Unlock() // want "mutex 'mu' has defer unlock but no corresponding lock"
		<-ch
	}()
}

func badRLockWithoutRUnlockCase() {
	var rw sync.RWMutex
	rw.RLock() // want "rwmutex 'rw' is rlocked but not runlocked"
}

func badRWUnlockCase() {
	var rw sync.RWMutex
	rw.RUnlock() // want "rwmutex 'rw' is runlocked but not rlocked"
}

func goodMutexDeferCase() {
	var mu sync.Mutex
	mu.Lock()
	defer mu.Unlock()
}

func goodConditionalBothBranchesCase() {
	var mu sync.Mutex
	cond := true
	if cond {
		mu.Lock()
		defer mu.Unlock()
	} else {
		mu.Lock()
		defer mu.Unlock()
	}
}

func goodRWMultipleOperationsCase() {
	var rw sync.RWMutex
	rw.RLock()
	rw.RUnlock()
	rw.Lock()
	rw.Unlock()
}
