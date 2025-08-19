# Gombinator <<VERSION>>

[![Go Report Card](https://goreportcard.com/badge/github.com/srgrr/gombinator)](https://goreportcard.com/report/github.com/srgrr/gombinator)

A goroutine-friendly functional library. It features methods like cartesian product for slices but by *generating* them on demand and channeling the results as you go.

# Quick Example

This code gets the first 10 natural numbers, filters the even numbers and computes their squares

```go
<<SAMPLE>>
```


## Some Notions

All functions require the user to provide a `context.Context` object. This allows the library to safely cancel results streaming prematurely.

There are two kinds of functions: normal functions and Cfunctions. Both **channel** their results and compute stuff **lazily**.

Cfunctions also work with **channels**. As you've seen in the sample, this means that you can chain different functions to perform complex lazy computations.
