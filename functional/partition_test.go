package functional

import (
	"context"
	"fmt"
	"testing"
)

func TestPartition(t *testing.T) {
	ctx := context.Background()
	strings := []string{"the", "quick", "brown", "fox", "jumps", "over", "the"}
	for partition := range Partition(ctx, strings, 2) {
		// TODO: WRITE AN ACTUAL TEST
		fmt.Printf("%v\n", partition)
	}
}
