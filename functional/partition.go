package functional

import "context"

// Partition channels consecutive array slices of size k
// except for maybe the last, which can be of size n % k
// Zero and negative numbers are not accepted here
func Partition[T any](ctx context.Context, A []T, k int) <-chan []T {
	if k < 1 {
		panic("k must be at least 1")
	}
	ch := make(chan []T)
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
	return ch
}
