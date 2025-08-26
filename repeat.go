package gombinator

import (
	"context"
	"errors"
)

var ErrInvalidRepeatKValue = errors.New("k must be -1 for infinite repetition, 0 for empty channel, or a positive integer")

// Repeat channels the same element k times
// It is similar to the built-in repeat function but works with channels
// If k is 0, it returns an empty channel
// If k is -1, it repeats indefinitely
// If k is positive, it repeats exactly k times
// Other values of k are not accepted
// It is similar to the built-in repeat function but works with channels
func Repeat[T any](ctx context.Context, elem T, k int) (<-chan T, error) {
	if k < -1 {
		return nil, ErrInvalidRepeatKValue
	}
	ch := make(chan T, GBufferSize)
	go func() {
		defer close(ch)
		for i := 0; k == -1 || i < k; i++ {
			select {
			case <-ctx.Done():
				return
			case ch <- elem:
			}
		}
	}()
	return ch, nil
}
