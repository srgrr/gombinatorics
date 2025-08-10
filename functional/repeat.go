package functional

func Repeat[T any](elem T, k int) chan T {
	ch := make(chan T)
	go func() {
		for i := 0; i < k; i++ {
			ch <- elem
		}
		close(ch)
	}()
	return ch
}
