package main

import (
	"context"
	"fmt"

	f "github.com/srgrr/gombinatorics/functional"
)

func main() {
	ctx := context.Background()
	evenSquaredNumbers :=
		f.CMap( // 3. Map the even numbers to their squares
			ctx,
			f.CFilter( // 2. Filter even numbers from the channel
				ctx,
				f.Range(ctx, 1, 11), // 1. Channel numbers from 1 to 10
				func(n int) bool { return n%2 == 0 },
			),
			func(n int) int { return n * n },
		)
	for n := range evenSquaredNumbers {
		fmt.Println(n)
	}
}
