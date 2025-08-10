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
