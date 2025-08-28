// Gombinator is a functional-style library designed to play nicely with goroutines and channels
// It allows the user to use functions like map, filter, partition, etc
// Most functions require the user to provide a context for early cancellation
// This library is NOT meant as a replacement for traditional loops: channels add overhead
// Consider using gombinator only when you were intending to distribute load between
// different goroutines in the first place
package gombinator

import "context"

// Buffer size for channels
var GBufferSize = 64

// Channelize channels the elements of a slice
func Channelize[T any](ctx context.Context, A []T) <-chan T {
	return Map(ctx, A, func(a T) T { return a })
}

// Filter channels elements from a slice that satisfy a given criterion
func Filter[T any](ctx context.Context, A []T, criterion func(T) bool) <-chan T {
	ch := make(chan T, GBufferSize)
	go func() {
		defer close(ch)
		for _, elem := range A {
			select {
			case <-ctx.Done():
				return
			default:
				if criterion(elem) {
					ch <- elem
				}
			}
		}
	}()
	return ch
}

// CFilter channels elements from a read-only channel that satisfy a given criterion
func CFilter[T any](ctx context.Context, A <-chan T, criterion func(T) bool) <-chan T {
	ch := make(chan T, GBufferSize)
	go func() {
		defer close(ch)
		for elem := range A {
			select {
			case <-ctx.Done():
				return
			default:
				if criterion(elem) {
					ch <- elem
				}
			}
		}
	}()
	return ch
}

// Map applies a given function to a given array and channels the elements
func Map[S any, T any](ctx context.Context, A []S, f func(S) T) <-chan T {
	ch := make(chan T, GBufferSize)
	go func() {
		defer close(ch)
		for _, elem := range A {
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

// CMap applies a function to a given read-only channel and channels the results
func CMap[S any, T any](ctx context.Context, A <-chan S, f func(S) T) <-chan T {
	ch := make(chan T, GBufferSize)
	go func() {
		defer close(ch)
		for elem := range A {
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

// Partition channels consecutive array slices of size k
// except for maybe the last, which can be of size n % k
// Zero and negative numbers are not accepted here
// Returns ErrInvalidKValue if k is not at least 1
func Partition[T any](ctx context.Context, A []T, k int) (<-chan []T, error) {
	if k < 1 {
		return nil, ErrInvalidPartitionKValue
	}
	ch := make(chan []T, GBufferSize)
	go func() {
		defer close(ch)
		for i := 0; i < len(A); i += k {
			select {
			case <-ctx.Done():
				return
			case ch <- A[i:min(len(A), i+k)]:
			}
		}
	}()
	return ch, nil
}

// CPartition channels consecutive slices of size k from a read-only channel
// except for maybe the last, which can be of size n % k
// Zero and negative numbers are not accepted here
// It is similar to Partition but works with channels instead of arrays
// Returns ErrInvalidKValue if k is not at least 1
func CPartition[T any](ctx context.Context, A <-chan T, k int) (<-chan []T, error) {
	if k < 1 {
		return nil, ErrInvalidPartitionKValue
	}
	ch := make(chan []T, GBufferSize)
	go func() {
		defer close(ch)
		batch := make([]T, 0, k)
		for elem := range A {
			select {
			case <-ctx.Done():
				return
			default:
				batch = append(batch, elem)
				if len(batch) == k {
					ch <- batch
					batch = make([]T, 0, k)
				}
			}
		}
		if len(batch) > 0 {
			ch <- batch
		}
	}()
	return ch, nil
}

// Range channels integers in [l, r), excluding r
// It is similar to the built-in range function but works with channels
func Range(ctx context.Context, l, r int) <-chan int {
	ch := make(chan int, GBufferSize)
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

// StepRange channels integers in [l, r) with a step
// It is similar to the built-in range function but works with channels
func StepRange(ctx context.Context, l, r, step int) <-chan int {
	if step <= 0 {
		panic("step must be greater than 0")
	}
	ch := make(chan int, GBufferSize)
	go func() {
		defer close(ch)
		for i := l; i < r; i += step {
			select {
			case <-ctx.Done():
				return
			case ch <- i:
			}
		}
	}()
	return ch
}

// Zips two slices and channels the corresponding pairs
// Zip won't fail if A or B are of different sizes, it'll
// just keep making pairs until one of the two slices runs
// out of elements
func Zip[P any, Q any](ctx context.Context, A []P, B []Q) <-chan Pair[P, Q] {
	ch := make(chan Pair[P, Q], GBufferSize)
	go func() {
		defer close(ch)
		for i := 0; i < min(len(A), len(B)); i++ {
			select {
			case <-ctx.Done():
				return
			default:
				ch <- Pair[P, Q]{First: A[i], Second: B[i]}
			}
		}
	}()
	return ch
}

// CZip zips two channels and channels the corresponding pairs
// Zip won't fail if A or B are of different sizes, it'll
// just keep making pairs until one of the two sources runs
// out of elements
func CZip[P any, Q any](ctx context.Context, A <-chan P, B <-chan Q) <-chan Pair[P, Q] {
	ch := make(chan Pair[P, Q], GBufferSize)
	go func() {
		defer close(ch)
		for {
			var a P
			var b Q
			var okA, okB bool

			select {
			case <-ctx.Done():
				return
			default:
				a, okA = <-A
				if !okA {
					return
				}
			}

			select {
			case <-ctx.Done():
				return
			default:
				b, okB = <-B
				if !okB {
					return
				}
			}

			select {
			case <-ctx.Done():
				return
			default:
				ch <- Pair[P, Q]{First: a, Second: b}
			}
		}
	}()
	return ch
}

// ERepeat channels the same element k times
// It is similar to the built-in repeat function but works with channels
// If k is 0, it returns an empty channel
// If k is -1, it repeats indefinitely
// If k is positive, it repeats exactly k times
// Other values of k are not accepted
// It is similar to the built-in repeat function but works with channels
func ERepeat[T any](ctx context.Context, elem T, k int) (<-chan T, error) {
	if k < -1 {
		return nil, ErrInvalidRepeatKValue
	}
	ch := make(chan T, GBufferSize)
	go func() {
		defer close(ch)
		for i := 0; k == -1 || i < k; i++ {
			select {
			case <-ctx.Done():
				return
			case ch <- elem:
			}
		}
	}()
	return ch, nil
}

// Repeat behaves like ERepeat but panics if k is not valid
func Repeat[T any](ctx context.Context, elem T, k int) <-chan T {
	rChan, err := ERepeat(ctx, elem, k)
	if err != nil {
		panic(err)
	}
	return rChan
}
