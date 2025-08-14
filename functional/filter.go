package functional

import "context"

// Filter channels elements from a slice that satisfy a given criterion
func Filter[T any](ctx context.Context, A []T, criterion func(T) bool) <-chan T {
	ch := make(chan T)
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
	ch := make(chan T)
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
