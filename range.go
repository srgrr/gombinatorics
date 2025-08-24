package gombinator

import "context"

// Range channels integers in [l, r), excluding r
// It is similar to the built-in range function but works with channels
func Range(ctx context.Context, l, r int) <-chan int {
	ch := make(chan int, GBufferSize)
	go func() {
		defer close(ch)
		for i := l; i < r; i++ {
			select {
			case <-ctx.Done():
				return
			case ch <- i:
			}
		}
	}()
	return ch
}

// StepRange channels integers in [l, r) with a step
// It is similar to the built-in range function but works with channels
func StepRange(ctx context.Context, l, r, step int) <-chan int {
	if step <= 0 {
		panic("step must be greater than 0")
	}
	ch := make(chan int, GBufferSize)
	go func() {
		defer close(ch)
		for i := l; i < r; i += step {
			select {
			case <-ctx.Done():
				return
			case ch <- i:
			}
		}
	}()
	return ch
}
