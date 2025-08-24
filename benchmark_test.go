package gombinator

import (
	"context"
	"testing"
)

var LIMIT = 10000000
var MOD int = 1e9 + 7

func BenchmarkRangeMapFilter_Gombinator(b *testing.B) {
	b.ReportAllocs()
	b.SetBytes(1)
	for n := 0; n < b.N; n++ {
		var total int = 0
		ctx := context.Background()
		squaredEvenNumbers :=
			CMap(ctx,
				CFilter(ctx,
					Range(ctx, 0, LIMIT),
					func(n int) bool { return n%2 == 0 },
				),
				func(n int) int { return n * n },
			)
		// Consume the channel for benchmarking
		for elem := range squaredEvenNumbers {
			total = (total + elem) % MOD
		}
		_ = total
	}
}

func getRangeChan(l, r int) <-chan int {
	ch := make(chan int, 64)
	go func() {
		defer close(ch)
		for i := l; i < r; i++ {
			ch <- i
		}
	}()
	return ch
}

func getMapChan(inputChan <-chan int, f func(n int) int) <-chan int {
	ch := make(chan int, 64)
	go func() {
		defer close(ch)
		for elem := range inputChan {
			ch <- f(elem)
		}
	}()
	return ch
}

func getFilterChan(inputChan <-chan int, f func(n int) bool) <-chan int {
	ch := make(chan int, 64)
	go func() {
		defer close(ch)
		for elem := range inputChan {
			if f(elem) {
				ch <- elem
			}
		}
	}()
	return ch
}

func BenchmarkRangeMapFilter_Vanilla(b *testing.B) {
	b.ReportAllocs()
	b.SetBytes(1)
	for n := 0; n < b.N; n++ {
		var total int = 0
		squaredEvenNumbers :=
			getMapChan(
				getFilterChan(
					getRangeChan(0, LIMIT),
					func(n int) bool { return n%2 == 0 },
				),
				func(n int) int { return n * n },
			)
		// Consume the channel for benchmarking
		for elem := range squaredEvenNumbers {
			total = (total + elem) % MOD
		}
		_ = total
	}
}

func Benchmark_Sequential(b *testing.B) {
	b.ReportAllocs()
	b.SetBytes(1)
	for n := 0; n < b.N; n++ {
		total := 0
		for i := 0; i < LIMIT; i++ {
			if i%2 == 0 {
				total = (total + i*i) % MOD
			}
		}
		_ = total
	}
}
