package gombinator

import (
	"context"
	"testing"
)

var LIMIT = 10000000
var MOD int = 1e9 + 7
var TARGET int = 3691500

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
		if total != TARGET {
			b.Errorf("Benchmark doesn't compute the right value")
		}
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
		if total != TARGET {
			b.Errorf("Benchmark doesn't compute the right value")
		}
	}
}

func getRangeChanCtx(ctx context.Context, l, r int) <-chan int {
	ch := make(chan int, 64)
	go func() {
		defer close(ch)
		for i := l; i < r; i++ {
			select {
			case <-ctx.Done():
				return
			default:
				ch <- i
			}
		}
	}()
	return ch
}

func getMapChanCtx(ctx context.Context, inputChan <-chan int, f func(n int) int) <-chan int {
	ch := make(chan int, 64)
	go func() {
		defer close(ch)
		for elem := range inputChan {
			select {
			case <-ctx.Done():
				return
			default:
				ch <- f(elem)
			}
		}
	}()
	return ch
}

func getFilterChanCtx(ctx context.Context, inputChan <-chan int, f func(n int) bool) <-chan int {
	ch := make(chan int, 64)
	go func() {
		defer close(ch)
		for elem := range inputChan {
			select {
			case <-ctx.Done():
				return
			default:
				if f(elem) {
					ch <- elem
				}
			}
		}
	}()
	return ch
}

func BenchmarkRangeMapFilter_VanillaContext(b *testing.B) {
	b.ReportAllocs()
	b.SetBytes(1)
	ctx := context.Background()
	for n := 0; n < b.N; n++ {
		var total int = 0
		squaredEvenNumbers :=
			getMapChanCtx(ctx,
				getFilterChanCtx(ctx,
					getRangeChanCtx(ctx, 0, LIMIT),
					func(n int) bool { return n%2 == 0 },
				),
				func(n int) int { return n * n },
			)
		// Consume the channel for benchmarking
		for elem := range squaredEvenNumbers {
			total = (total + elem) % MOD
		}
		if total != TARGET {
			b.Errorf("Benchmark doesn't compute the right value")
		}
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
		if total != TARGET {
			b.Errorf("Benchmark doesn't compute the right value")
		}
	}
}
