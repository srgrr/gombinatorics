package gombinator

import "context"

func Channelize[T any](ctx context.Context, A []T) <-chan T {
	return Map(ctx, A, func(a T) T { return a })
}
