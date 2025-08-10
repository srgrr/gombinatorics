package sets

import "context"

func subsets[T any](ctx context.Context, A []T, subset []T, limit int, ch chan []T) {
	select {
	case <-ctx.Done():
		return
	default:
		if len(A) == 0 {
			if limit == -1 || limit == 0 {
				select {
				case <-ctx.Done():
					return
				case ch <- subset:
				}
			}
			return
		}
		current := A[len(A)-1]
		A = A[:len(A)-1]
		subsets(ctx, A, subset, limit, ch)
		if limit == -1 || limit > 0 {
			subset = append(subset, current)
			subsets(ctx, A, subset, max(-1, limit-1), ch)
		}
	}
}
