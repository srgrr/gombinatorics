package gombinator

import "context"

// Channelize channels the elements of a slice
func Channelize[T any](ctx context.Context, A []T) <-chan T {
	return Map(ctx, A, func(a T) T { return a })
}
