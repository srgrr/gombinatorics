# Gombinatorics 🎲

[![Go Report Card](https://goreportcard.com/badge/github.com/srgrr/gombinatorics)](https://goreportcard.com/report/github.com/srgrr/gombinatorics)

A goroutine-friendly combinatorics/functional library. It features methods like cartesian product for slices but by *generating* them on demand and channeling the results as you go.

## ⚠️ Friendly Warning ⚠️
This library is still WIP. It started as a side-quest for something I'm working on.

I'll keep adding samples, configurations and support as time goes by.

# Functional example
The library allows you to turn code like this
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

Into code like this

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