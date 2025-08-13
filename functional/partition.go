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

// CPartition channels consecutive slices of size k from a read-only channel
// except for maybe the last, which can be of size n % k
// Zero and negative numbers are not accepted here
// It is similar to Partition but works with channels instead of arrays
func CPartition[T any](ctx context.Context, A <-chan T, k int) <-chan []T {
	if k < 1 {
		panic("k must be at least 1")
	}
	ch := make(chan []T)
	go func() {
		defer close(ch)
		batch := make([]T, 0, k)
		for elem := range A {
			select {
			case <-ctx.Done():
				return
			default:
				batch = append(batch, elem)
				if len(batch) == k {
					ch <- batch
					batch = make([]T, 0, k)
				}
			}
		}
		if len(batch) > 0 {
			ch <- batch
		}
	}()
	return ch
}
