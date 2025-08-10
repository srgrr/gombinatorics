package functional

import "context"

// Map applies a given function to a given array and channels the elements
func Map[S any, T any](ctx context.Context, A []S, f func(S) T) <-chan T {
	ch := make(chan T)
	go func() {
		defer close(ch)
		for _, elem := range A {
			select {
			case <-ctx.Done():
				return
			case ch <- f(elem):
			}
		}
	}()
	return ch
}

// CMap applies a function to a given read-only channel and channels the results
func CMap[S any, T any](ctx context.Context, A <-chan S, f func(S) T) <-chan T {
	ch := make(chan T)
	go func() {
		defer close(ch)
		for elem := range A {
			select {
			case <-ctx.Done():
				return
			case ch <- f(elem):
			}
		}
	}()
	return ch
}
