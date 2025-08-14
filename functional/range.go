package functional

import "context"

// Range channels integers in [l, r), excluding r
// It is similar to the built-in range function but works with channels
func Range(ctx context.Context, l, r int) <-chan int {
	ch := make(chan int)
	go func() {
		defer close(ch)
		for i := l; i < r; i++ {
			select {
			case <-ctx.Done():
				return
			case ch <- i:
			}
		}
	}()
	return ch
}
