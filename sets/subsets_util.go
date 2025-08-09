package sets

func subsets[T any](A []T, subset []T, limit int, ch chan []T) {
	if len(A) == 0 {
		if limit == -1 || limit == 0 {
			ch <- subset
		}
		return
	}
	current := A[len(A)-1]
	A = A[:len(A)-1]
	subsets(A, subset, limit, ch)
	if limit == -1 || limit > 0 {
		subset = append(subset, current)
		subsets(A, subset, max(-1, limit-1), ch)
	}
}
