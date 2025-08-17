package sets

import (
	"context"

	types "github.com/srgrr/gombinatorics/types"
)

// Generates the cartesian product of two slices
// Slices can be of different types
// Pairs are wrapped in a types.Pair
// Elements are paired following the given order in both arrays, from left to right
func CartesianProduct[P any, Q any](ctx context.Context, A []P, B []Q) <-chan types.Pair[P, Q] {
	ch := make(chan types.Pair[P, Q])
	go func() {
		defer close(ch)
		for _, a := range A {
			for _, b := range B {
				select {
				case <-ctx.Done():
					return
				case ch <- types.Pair[P, Q]{First: a, Second: b}:
				}
			}
		}
	}()
	return ch
}
