# Gombinatorics 🎲 & 🐑 (da)

[![Go Report Card](https://goreportcard.com/badge/github.com/srgrr/gombinatorics)](https://goreportcard.com/report/github.com/srgrr/gombinatorics)

A goroutine-friendly combinatorics/functional library. It features methods like cartesian product for slices but by *generating* them on demand and channeling the results as you go.

## Quick Overview
### Functional
`functional` provides the most usual functional programming patterns but adapted to channels. That is, they're designed to **channel** the results instead of accumulating them and returning whole computed collections. The advantages of this are twofold: all elements are computed on demand and can be consumed by different goroutines with no concurrency issues.

`functional` provides the following functions:
```go
package functional // import "github.com/srgrr/gombinatorics/functional"


FUNCTIONS

**func CFilter[T any](ctx context.Context, A <-chan T, criterion func(T) bool) <-chan T**
```go
CFilter channels elements from a read-only channel that satisfy a given
```
```go
criterion
```

**func CMap[S any, T any](ctx context.Context, A <-chan S, f func(S) T) <-chan T**
```go
CMap applies a function to a given read-only channel and channels the
```
```go
results
```

**func CPartition[T any](ctx context.Context, A <-chan T, k int) <-chan []T**
```go
CPartition channels consecutive slices of size k from a read-only channel
```
```go
except for maybe the last, which can be of size n % k Zero and negative
```
```go
numbers are not accepted here It is similar to Partition but works with
```
```go
channels instead of arrays
```

**func CZip[P any, Q any](ctx context.Context, A <-chan P, B <-chan Q) <-chan types.Pair[P, Q]**
```go
CZip zips two channels and channels the corresponding pairs Zip won't fail
```
```go
if A or B are of different sizes, it'll just keep making pairs until one of
```
```go
the two sources runs out of elements
```

**func Filter[T any](ctx context.Context, A []T, criterion func(T) bool) <-chan T**
```go
Filter channels elements from a slice that satisfy a given criterion
```

**func Map[S any, T any](ctx context.Context, A []S, f func(S) T) <-chan T**
```go
Map applies a given function to a given array and channels the elements
```

**func Partition[T any](ctx context.Context, A []T, k int) <-chan []T**
```go
Partition channels consecutive array slices of size k except for maybe the
```
```go
last, which can be of size n % k Zero and negative numbers are not accepted
```
```go
here
```

**func Range(ctx context.Context, l, r int) <-chan int**
```go
Range channels integers in [l, r), excluding r It is similar to the built-in
```
```go
range function but works with channels
```

**func Repeat[T any](ctx context.Context, elem T, k int) chan T**
**func Zip[P any, Q any](ctx context.Context, A []P, B []Q) <-chan types.Pair[P, Q]**
```go
Zips two slices and channels the corresponding pairs Zip won't fail if A or
```
```go
B are of different sizes, it'll just keep making pairs until one of the two
```
```go
slices runs out of elements
```

```
### Sets
`sets` provides the following functions:
```go
package sets // import "github.com/srgrr/gombinatorics/sets"


FUNCTIONS

**func CartesianProduct[P any, Q any](ctx context.Context, A []P, B []Q) <-chan types.Pair[P, Q]**
```go
Generates the cartesian product of two slices Slices can be of different
```
```go
types Pairs are wrapped in a types.Pair Elements are paired following the
```
```go
given order in both arrays, from left to right
```

**func Subsets[T any](ctx context.Context, A []T) <-chan []T**
```go
Returns all the subsets Elements are included (excluded) in the given order
```

**func SubsetsOfFixedSize[T any](ctx context.Context, A []T, k int) <-chan []T**
```go
Returns all the subsets whose size is exactly k
```

```

# Functional example
The library allows you to turn memory-heavy, single threaded code like this
```go
func main() {
	// 1. Create a slice of numbers from 1 to 10
	numbers := make([]int, 10)
	for i := 1; i <= 10; i++ {
		numbers[i-1] = i
	}
	// 2. Filter even numbers
	evenNumbers := make([]int, 0)
	for _, n := range numbers {
		if n%2 == 0 {
			evenNumbers = append(evenNumbers, n)
		}
	}
	// 3. Map the even numbers to their squares
	squaredEvenNumbers := make([]int, len(evenNumbers))
	for i, n := range evenNumbers {
		squaredEvenNumbers[i] = n * n
	}
	for _, n := range squaredEvenNumbers {
		println(n)
	}
}
```

Into lightweight, *on-demand* code like this

```go
package main

import (
	"context"
	"fmt"

	f "github.com/srgrr/gombinatorics/functional"
)

func main() {
	ctx := context.Background()
	evenSquaredNumbers :=
		f.CMap( // 3. Map the even numbers to their squares
			ctx,
			f.CFilter( // 2. Filter even numbers from the channel
				ctx,
				f.Range(ctx, 1, 11), // 1. Channel numbers from 1 to 10
				func(n int) bool { return n%2 == 0 },
			),
			func(n int) int { return n * n },
		)
	for n := range evenSquaredNumbers {
		fmt.Println(n)
	}
}
```

Functions starting with `C` **channel** their results. That is, this second sample computes elements **on demand** instead of computing whole lists.
Furthermore, `evenSquaredNumbers` can be consumed by different goroutines at the same time!
