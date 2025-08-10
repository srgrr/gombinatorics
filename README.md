# Gombinatorics 🎲

A goroutine-friendly combinatorics/functional library. It features methods like cartesian product for slices but by *generating* them on demand and channeling the results as you go.

## ⚠️ Friendly Warning ⚠️
This library is still WIP. It started as a side-quest for something I'm working on.

Immediate TODOs:
- Add `CFunctions` accepting channels, so they can also take *generators*

# Example
This example prints `hello world`.
```go
package main

import (
	"fmt"

	f "github.com/srgrr/gombinatorics/functional"
	t "github.com/srgrr/gombinatorics/types"
)

func main() {
	prefixes := []string{"hel", "wor"}
	suffixes := []string{"lo", "ld"}
	concat := func(p t.Pair[string, string]) string {
		return p.First + p.Second
	}
	for word := range f.CMap(f.Zip(prefixes, suffixes), concat) {
		fmt.Printf("%s ", word)
	}
	fmt.Println()
}
```