package gombinatorics

func CartesianProduct[P any, Q any](A []P, B []Q) <-chan Pair[P, Q] {
	ch := make(chan Pair[P, Q])
	go func() {
		for _, a := range A {
			for _, b := range B {
				ch <- Pair[P, Q]{a, b}
			}
		}
		close(ch)
	}()
	return ch
}

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

func Subsets[T any](A []T) <-chan []T {
	ch := make(chan []T)
	go func() {
		subset := make([]T, 0)
		subsets(A, subset, -1, ch)
		close(ch)
	}()
	return ch
}

func SubsetsOfFixedSize[T any](A []T, k int) <-chan []T {
	ch := make(chan []T)
	go func() {
		subset := make([]T, 0)
		subsets(A, subset, k, ch)
		close(ch)
	}()
	return ch
}

func Zip[P any, Q any](A []P, B []Q) <-chan Pair[P, Q] {
	ch := make(chan Pair[P, Q])
	go func() {
		for i := 0; i < min(len(A), len(B)); i++ {
			ch <- Pair[P, Q]{A[i], B[i]}
		}
		close(ch)
	}()
	return ch
}
