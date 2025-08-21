# Gombinator v1.0.3

[![Go Report Card](https://goreportcard.com/badge/github.com/srgrr/gombinator)](https://goreportcard.com/report/github.com/srgrr/gombinator)

A goroutine-friendly functional library. It features methods like `map`, `filter`, etc for slices and channels but by *generating* the results on demand.

# How to use it
This library can be installed and imported as a dependency by doing the usually `go get` stuff...
```
go get github.com/srgrr/gombinator
```
... and import it the same way
```go
import g "github.com/srgrr/gombinator"
```

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

# Declarative Example

You can declare *what* you're intending to do and do it afterwards

```go
// Find the whole example in examples/declarative_example
data := []string{"hello", "world"}

// Declare channeling functions
lenChan := g.Map(ctx, data, func(s string) int { return len(s) })
filterChan := g.Filter(ctx, data, func(s string) bool { return s == "hello" })

// Run them!
go lenTask(lenChan)
go helloTask(filterChan)

```

# Some Notions

All functions require the user to provide a `context.Context` object. This allows the library to safely cancel results streaming prematurely.

There are two kinds of functions: normal functions and Cfunctions. Both **channel** their results and compute stuff **lazily**.

Cfunctions also work with **channels**. As you've seen in the sample, this means that you can chain different functions to perform complex lazy computations.

# Things to Watch Out For

The library does **more** than just declare stuff: it opens actual channels and leaves them blocked waiting for someone to read from them.

This means two things:
    - Declaring stuff **does** add overhead, and so does using channels in general. This lib is meant to ease writing pipelines meant to be consumed by multiple goroutines, not to write fancy pythonic oneliners
    - You're responsible of avoiding deadlocks accordingly

