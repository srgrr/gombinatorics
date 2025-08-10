package functional

import (
	"context"

	types "github.com/srgrr/gombinatorics/types"
)

// Zips two slices and channels the corresponding pairs
// Zip won't fail if A or B are of different sizes, it'll
// just keep making pairs until one of the two slices runs
// out of elements
func Zip[P any, Q any](ctx context.Context, A []P, B []Q) <-chan types.Pair[P, Q] {
	ch := make(chan types.Pair[P, Q])
	go func() {
		defer close(ch)
		for i := 0; i < min(len(A), len(B)); i++ {
			select {
			case <-ctx.Done():
				return
			case ch <- types.Pair[P, Q]{First: A[i], Second: B[i]}:
			}
		}
	}()
	return ch
}

// CZip zips two channels and channels the corresponding pairs
// Zip won't fail if A or B are of different sizes, it'll
// just keep making pairs until one of the two sources runs
// out of elements
func CZip[P any, Q any](ctx context.Context, A <-chan P, B <-chan Q) <-chan types.Pair[P, Q] {
	ch := make(chan types.Pair[P, Q])
	go func() {
		defer close(ch)
		for {
			var a P
			var b Q
			var okA, okB bool

			select {
			case <-ctx.Done():
				return
			case a, okA = <-A:
				if !okA {
					return
				}
			}

			select {
			case <-ctx.Done():
				return
			case b, okB = <-B:
				if !okB {
					return
				}
			}

			select {
			case <-ctx.Done():
				return
			case ch <- types.Pair[P, Q]{First: a, Second: b}:
			}
		}
	}()
	return ch
}
