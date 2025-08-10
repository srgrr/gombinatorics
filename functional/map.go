package functional

func Map[S any, T any](A []S, f func(S) T) <-chan T {
	ch := make(chan T)
	go func() {
		for _, elem := range A {
			ch <- f(elem)
		}
		close(ch)
	}()
	return ch
}

func CMap[S any, T any](A <-chan S, f func(S) T) <-chan T {
	ch := make(chan T)
	go func() {
		for elem := range A {
			ch <- f(elem)
		}
		close(ch)
	}()
	return ch
}
