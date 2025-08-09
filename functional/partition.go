package functional

// Partition channels consecutive array slices of size k
// except for maybe the last, which can be of size n % k
// Zero and negative numbers are not accepted here
func Partition[T any](A []T, k int) <-chan []T {
	if k < 1 {
		panic("k must be at least 1")
	}
	ch := make(chan []T)
	go func() {
		for i := 0; i < len(A); i += k {
			ch <- A[i:min(len(A), i+k)]
		}
		close(ch)
	}()
	return ch
}
