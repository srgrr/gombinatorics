# Gombinator v1.0.2

[![Go Report Card](https://goreportcard.com/badge/github.com/srgrr/gombinator)](https://goreportcard.com/report/github.com/srgrr/gombinator)

A goroutine-friendly functional library. It features methods like cartesian product for slices but by *generating* them on demand and channeling the results as you go.

# Quick Example

This code gets the first 10 natural numbers, filters the even numbers and computes their squares

```go
package main

import (
	"context"
	"fmt"

	g "github.com/srgrr/gombinator"
)

func main() {
	ctx := context.Background()
	evenSquaredNumbers :=
		g.CMap(ctx, // 3. Map the even numbers to their squares
			g.CFilter(ctx, // 2. Filter even numbers from the channel
				g.Range(ctx, 1, 11), // 1. Channel numbers from 1 to 10
				func(n int) bool { return n%2 == 0 },
			),
			func(n int) int { return n * n },
		)
	for n := range evenSquaredNumbers {
		fmt.Println(n)
	}
}

```


## Some Notions

All functions require the user to provide a `context.Context` object. This allows the library to safely cancel results streaming prematurely.

There are two kinds of functions: normal functions and Cfunctions. Both **channel** their results and compute stuff **lazily**.

Cfunctions also work with **channels**. As you've seen in the sample, this means that you can chain different functions to perform complex lazy computations.
