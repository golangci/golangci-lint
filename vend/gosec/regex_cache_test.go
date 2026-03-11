package gosec

import (
	"fmt"
	"regexp"
	"sync"
	"testing"
)

func TestGlobalCache_Stress(t *testing.T) {
	// Simple stress test to ensure thread safety (running with -race is ideal)
	// We can't easily assert on race conditions without the race detector,
	// but this ensures no obvious panics or deadlocks.

	const routines = 10
	const iterations = 100

	// Use a test regex for the cache key
	testRe := regexp.MustCompile(`test`)

	var wg sync.WaitGroup
	wg.Add(routines)

	for i := range routines {
		go func(id int) {
			defer wg.Done()
			key := regexCacheKey{Re: testRe, Str: fmt.Sprintf("str-%d", id)}
			for j := range iterations {
				GlobalCache.Add(key, j)
				if _, ok := GlobalCache.Get(key); !ok {
					t.Errorf("failed to get key %v", key)
				}
			}
		}(i)
	}
	wg.Wait()
}
