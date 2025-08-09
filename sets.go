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

func subsets[T any](A []T, subset []T, ch chan []T) {
	if len(A) == 0 {
		ch <- subset
		return
	}
	current := A[len(A)-1]
	A = A[:len(A)-1]
	subsets(A, subset, ch)
	subset = append(subset, current)
	subsets(A, subset, ch)
}

func Subsets[T any](A []T) <-chan []T {
	ch := make(chan []T)
	go func() {
		subset := make([]T, 0)
		subsets(A, subset, ch)
		close(ch)
	}()
	return ch
}
