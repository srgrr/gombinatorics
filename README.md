# Gombinator v1.1.0

[![Go Report Card](https://goreportcard.com/badge/github.com/srgrr/gombinator)](https://goreportcard.com/report/github.com/srgrr/gombinator) [![Go Reference](https://pkg.go.dev/badge/github.com/srgrr/gombinator.svg)](https://pkg.go.dev/github.com/srgrr/gombinator) ![License](https://img.shields.io/github/license/srgrr/gombinator)




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

`gombinator` also allows to tune channel buffering (only for channels within the library itself) via the parameter `GBufferSize`.

# Things to Watch Out For

The library does **more** than just declaring stuff: it opens actual channels and leaves them blocked waiting for someone to read from them.

This means two things:

- Declaring stuff **does** add overhead, and so does using channels in general. This lib is meant to ease writing pipelines meant to be consumed by multiple goroutines, not to write fancy pythonic oneliners

- You're responsible of avoiding deadlocks. You still gotta manage context cancelling and channel consumption accordingly

# Performance

As we mentioned before, `gombinator` does take a toll on performance due to the heavy use of channels. You can check it out by yourself by running `benchmark_test.go`. Here are some results:

```
goos: darwin
goarch: arm64
pkg: github.com/srgrr/gombinator
cpu: Apple M4
BenchmarkRangeMapFilter_Gombinator-10                  2         919172438 ns/op           0.00 MB/s           2768 B/op          7 allocs/op
BenchmarkRangeMapFilter_Vanilla-10                     2         828944979 ns/op           0.00 MB/s           2608 B/op          8 allocs/op
BenchmarkRangeMapFilter_VanillaContext-10              2         924227458 ns/op           0.00 MB/s           2064 B/op          6 allocs/op
Benchmark_Sequential-10                              331           3581338 ns/op           0.00 MB/s              0 B/op          0 allocs/op
```

`-Gombinator`, `-Vanilla`, and `-VanillaContext` use channels. Both `-Vanilla` approaches implement the same channel-based approach as `gombinator` (and they all implement the same example where we generate some numbers, filter the even ones and compute their squares afterwards). `-Vanilla` is not context aware whereas `-VanillaContext` is. In fact, `-VanillaContext` can be considered an inline version of what `gombinator` does.

`-Sequential` implements the same thing but running a full sequential loop.

The goal of this benchmark is showing the user that `gombinator` is NOT a replacement for simple sequential code where it is enough doing so and that context awareness always comes with a cost.

On the bright side, `gombinator` allows you to save lots of code compared to having to write all the context handling and function pipelining by yourself.