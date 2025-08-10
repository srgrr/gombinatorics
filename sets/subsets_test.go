package sets

import (
	"context"
	"testing"
)

func TestSubsets(t *testing.T) {
	ctx := context.Background()
	nums := []int{1, 2, 4, 8}
	total := 0
	numSubsets := 0
	for subset := range Subsets(ctx, nums) {
		numSubsets++
		for _, elem := range subset {
			total += elem
		}
	}
	if total != 8*(1+2+4+8) || numSubsets != 16 {
		t.Errorf("%d %d", total, numSubsets)
	}
}

func TestSubsetsOfFixedSize(t *testing.T) {
	ctx := context.Background()
	nums := []int{1, 2, 4, 8}
	total := 0
	numSubsets := 0
	for subset := range SubsetsOfFixedSize(ctx, nums, 2) {
		numSubsets++
		for _, elem := range subset {
			total += elem
		}
	}
	if total != 3*(1+2+4+8) || numSubsets != 6 {
		t.Errorf("%d %d", total, numSubsets)
	}
}
