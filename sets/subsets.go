package sets

import "context"

// Subsets generates all subsets of a given slice A
// Returns all the subsets
// Elements are included (excluded) in the given order
func Subsets[T any](ctx context.Context, A []T) <-chan []T {
	ch := make(chan []T)
	go func() {
		defer close(ch)
		subset := make([]T, 0)
		subsets(ctx, A, subset, -1, ch)
	}()
	return ch
}

// SubsetsOfFixedSize generates all subsets of a given slice A
// Returns all the subsets whose size is exactly k
func SubsetsOfFixedSize[T any](ctx context.Context, A []T, k int) <-chan []T {
	ch := make(chan []T)
	go func() {
		defer close(ch)
		subset := make([]T, 0)
		subsets(ctx, A, subset, k, ch)
	}()
	return ch
}
