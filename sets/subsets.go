package sets

// Returns all the subsets
// Elements are included (excluded) in the given order
func Subsets[T any](A []T) <-chan []T {
	ch := make(chan []T)
	go func() {
		subset := make([]T, 0)
		subsets(A, subset, -1, ch)
		close(ch)
	}()
	return ch
}

// Returns all the subsets whose size is exactly k
func SubsetsOfFixedSize[T any](A []T, k int) <-chan []T {
	ch := make(chan []T)
	go func() {
		subset := make([]T, 0)
		subsets(A, subset, k, ch)
		close(ch)
	}()
	return ch
}
