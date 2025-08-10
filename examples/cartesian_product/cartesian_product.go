package main

import (
	"context"
	"fmt"

	sets "github.com/srgrr/gombinatorics/sets"
)

func main() {
	ctx := context.Background()
	emails := []string{"sergio@raccoon.me", "raquel@cat.me"}
	spam := []string{"Raccoon plushies now 10% discount", "Brown hair dye now 5% discount"}
	for emailAndSpam := range sets.CartesianProduct(ctx, emails, spam) {
		fmt.Printf(
			"%s got spammed with \"%s\"\n",
			emailAndSpam.First, emailAndSpam.Second,
		)
	}
}
