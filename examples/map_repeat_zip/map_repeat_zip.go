package main

import (
	"context"
	"fmt"

	f "github.com/srgrr/gombinator"
)

func main() {
	ctx := context.Background()
	prefixes := []string{"hel", "wor"}
	suffixes := []string{"lo", "ld"}
	concat := func(p f.Pair[string, string]) string {
		return p.First + p.Second
	}
	for word := range f.CMap(ctx, f.Zip(ctx, prefixes, suffixes), concat) {
		fmt.Printf("%s ", word)
	}
	fmt.Println()
}
