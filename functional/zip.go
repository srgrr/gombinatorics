package functional

import types "github.com/srgrr/gombinatorics/types"

func Zip[P any, Q any](A []P, B []Q) <-chan types.Pair[P, Q] {
	ch := make(chan types.Pair[P, Q])
	go func() {
		for i := 0; i < min(len(A), len(B)); i++ {
			ch <- types.Pair[P, Q]{First: A[i], Second: B[i]}
		}
		close(ch)
	}()
	return ch
}
