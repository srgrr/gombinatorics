package functional

import "context"

func Repeat[T any](ctx context.Context, elem T, k int) chan T {
	ch := make(chan T)
	go func() {
		defer close(ch)
		for i := 0; i < k; i++ {
			select {
			case <-ctx.Done():
				return
			case ch <- elem:
			}
		}
	}()
	return ch
}
