# Gombinatorics 🎲

A goroutine-friendly combinatorics/functional library. It features methods like cartesian product for slices but by *generating* them on demand and channeling the results as you go.

## ⚠️ Friendly Warning ⚠️
This library is still WIP. It started as a side-quest for something I'm working on.

Immediate TODOs:
- Add `IFunctions` accepting channels, so they can also take *generators*

# Example

```go
package main

import (
	"fmt"

	sets "github.com/srgrr/gombinatorics/sets"
)

func main() {
	emails := []string{"sergio@raccoon.me", "raquel@cat.me"}
	spam := []string{"Raccoon plushies now 10% discount", "Brown hair dye now 5% discount"}
	for emailAndSpam := range sets.CartesianProduct(emails, spam) {
		fmt.Printf(
			"%s got spammed with \"%s\"\n",
			emailAndSpam.First, emailAndSpam.Second,
		)
	}
}
```