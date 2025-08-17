package functional

import "context"

// Repeat channels the same element k times
// It is similar to the built-in repeat function but works with channels
// Zero and negative numbers are not accepted here
// If k is 0, it returns an empty channel
func Repeat[T any](ctx context.Context, elem T, k int) chan T {
	if k < 0 {
		panic("k must be at least 0")
	}
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
