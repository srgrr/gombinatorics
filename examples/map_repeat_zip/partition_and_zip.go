package main

import (
	"context"
	"fmt"

	f "github.com/srgrr/gombinatorics/functional"
	t "github.com/srgrr/gombinatorics/types"
)

func main() {
	ctx := context.Background()
	prefixes := []string{"hel", "wor"}
	suffixes := []string{"lo", "ld"}
	concat := func(p t.Pair[string, string]) string {
		return p.First + p.Second
	}
	for word := range f.CMap(ctx, f.Zip(ctx, prefixes, suffixes), concat) {
		fmt.Printf("%s ", word)
	}
	fmt.Println()
}
