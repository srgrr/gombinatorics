package sets

import types "github.com/srgrr/gombinatorics/types"

// Generates the cartesian product of two slices
// Slices can be of different types
// Pairs are wrapped in a types.Pair
// Elements are paired following the given order in both arrays, from left to right
func CartesianProduct[P any, Q any](A []P, B []Q) <-chan types.Pair[P, Q] {
	ch := make(chan types.Pair[P, Q])
	go func() {
		for _, a := range A {
			for _, b := range B {
				ch <- types.Pair[P, Q]{First: a, Second: b}
			}
		}
		close(ch)
	}()
	return ch
}
