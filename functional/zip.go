package functional

import types "github.com/srgrr/gombinatorics/types"

// Zips two slices and channels the corresponding pairs
// Zip won't fail if A or B are of different sizes, it'll
// just keep making pairs until one of the two slices runs
// out of elements
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

// Zips two channels and channels the corresponding pairs
// Zip won't fail if A or B are of different sizes, it'll
// just keep making pairs until one of the two sources runs
// out of elements
func CZip[P any, Q any](A <-chan P, B <-chan Q) <-chan types.Pair[P, Q] {
	ch := make(chan types.Pair[P, Q])
	go func() {
		ok := true
		for ok {
			a, ok := <-A
			var b Q
			if ok {
				b, ok = <-B
			}
			if ok {
				ch <- types.Pair[P, Q]{First: a, Second: b}
			}
		}
		close(ch)
	}()
	return ch
}
