//go:build go1.25

//golangcitest:args -Emodernize
//golangcitest:expected_exitcode 0
package bloop

import (
	"sync"
	"testing"
)

func BenchmarkA(b *testing.B) {
	println("slow")
	b.ResetTimer()

	for range b.N { // want "b.N can be modernized using b.Loop.."
	}
}

func BenchmarkB(b *testing.B) {
	// setup
	{
		b.StopTimer()
		println("slow")
		b.StartTimer()
	}

	for i := range b.N { // Nope. Should we change this to "for i := 0; b.Loop(); i++"?
		print(i)
	}

	b.StopTimer()
	println("slow")
}

func BenchmarkC(b *testing.B) {
	// setup
	{
		b.StopTimer()
		println("slow")
		b.StartTimer()
	}

	for i := 0; i < b.N; i++ { // want "b.N can be modernized using b.Loop.."
		println("no uses of i")
	}

	b.StopTimer()
	println("slow")
}

func BenchmarkD(b *testing.B) {
	for i := 0; i < b.N; i++ { // want "b.N can be modernized using b.Loop.."
		println(i)
	}
}

func BenchmarkE(b *testing.B) {
	b.Run("sub", func(b *testing.B) {
		b.StopTimer() // not deleted
		println("slow")
		b.StartTimer() // not deleted

		// ...
	})
	b.ResetTimer()

	for i := 0; i < b.N; i++ { // want "b.N can be modernized using b.Loop.."
		println("no uses of i")
	}

	b.StopTimer()
	println("slow")
}

func BenchmarkF(b *testing.B) {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < b.N; i++ { // nope: b.N accessed from a FuncLit
		}
	}()
	wg.Wait()
}

func BenchmarkG(b *testing.B) {
	var wg sync.WaitGroup
	poster := func() {
		for i := 0; i < b.N; i++ { // nope: b.N accessed from a FuncLit
		}
		wg.Done()
	}
	wg.Add(2)
	for i := 0; i < 2; i++ {
		go poster()
	}
	wg.Wait()
}

func BenchmarkH(b *testing.B) {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for range b.N { // nope: b.N accessed from a FuncLit
		}
	}()
	wg.Wait()
}

func BenchmarkI(b *testing.B) {
	for i := 0; i < b.N; i++ { // nope: b.N accessed more than once in benchmark
	}
	for i := 0; i < b.N; i++ { // nope: b.N accessed more than once in benchmark
	}
}
