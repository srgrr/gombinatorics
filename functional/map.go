package functional

import "context"

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
