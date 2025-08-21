# Gombinator <<VERSION>>

[![Go Report Card](https://goreportcard.com/badge/github.com/srgrr/gombinator)](https://goreportcard.com/report/github.com/srgrr/gombinator)

A goroutine-friendly functional library. It features methods like `map`, `filter`, etc for slices and channels but by *generating* the results on demand.

# Quick Example

This code gets the first 10 natural numbers, filters the even numbers and computes their squares

```go
<<SAMPLE>>
```

# Declarative Example

You can declare *what* you're intending to do and do it afterwards

```go
<<DECLARATIVE>>
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
