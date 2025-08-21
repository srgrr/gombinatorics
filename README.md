# Gombinator v1.0.2

[![Go Report Card](https://goreportcard.com/badge/github.com/srgrr/gombinator)](https://goreportcard.com/report/github.com/srgrr/gombinator)

A goroutine-friendly functional library. It features methods like `map`, `filter`, etc for slices and channels but by *generating* the results on demand.

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
package main

import (
	"context"
	"fmt"
	"sync"

	g "github.com/srgrr/gombinator"
)

func lenTask(wg *sync.WaitGroup, numChan <-chan int) {
	defer wg.Done()
	for numVal := range numChan {
		fmt.Println(numVal)
	}
}

func helloTask(wg *sync.WaitGroup, stringChan <-chan string) {
	defer wg.Done()
	for stringVal := range stringChan {
		fmt.Println(stringVal)
	}
}

func main() {
	ctx := context.Background()
	data := []string{"hello", "world"}
	// Declare channeling functions
	lenChan := g.Map(ctx, data, func(s string) int { return len(s) })
	filterChan := g.Filter(ctx, data, func(s string) bool { return s == "hello" })

	// Run them!
	var wg sync.WaitGroup
	wg.Add(2)
	go lenTask(&wg, lenChan)
	go helloTask(&wg, filterChan)
	wg.Wait()
}

```

# Important Caveats

The library does **more** than just declare stuff: it opens actual channels and leaves them blocked waiting for someone to read from them.

This means two things:
    - Declaring stuff **does** add overhead
    - You're responsible of avoiding deadlocks accordingly

## Some Notions

All functions require the user to provide a `context.Context` object. This allows the library to safely cancel results streaming prematurely.

There are two kinds of functions: normal functions and Cfunctions. Both **channel** their results and compute stuff **lazily**.

Cfunctions also work with **channels**. As you've seen in the sample, this means that you can chain different functions to perform complex lazy computations.
