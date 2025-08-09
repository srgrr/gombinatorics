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
