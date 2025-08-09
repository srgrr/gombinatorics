package functional

import (
	"fmt"
	"testing"
)

func TestPartition(t *testing.T) {
	strings := []string{"the", "quick", "brown", "fox", "jumps", "over", "the"}
	for partition := range Partition(strings, 2) {
		// TODO: WRITE AN ACTUAL TEST
		fmt.Printf("%v\n", partition)
	}
}
