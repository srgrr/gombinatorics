package gombinator

import (
	"context"
	"reflect"
	"testing"
)

type MapTestCase[T any] struct {
	name     string
	A        []T
	mapFunc  func(T) T
	expected []T
}

func TestMap(t *testing.T) {
	ctx := context.Background()
	tests := []MapTestCase[int]{
		{
			"Square each element",
			[]int{1, 2, 3, 4, 5},
			func(n int) int { return n * n },
			[]int{1, 4, 9, 16, 25},
		},
		{
			"Return empty list when input is empty",
			[]int{},
			func(n int) int { return n * n },
			[]int{},
		},
	}

	for _, tc := range tests {
		t.Run(
			tc.name,
			func(t *testing.T) {
				got := make([]int, 0)
				for elem := range Map(ctx, tc.A, tc.mapFunc) {
					got = append(got, elem)
				}
				if !reflect.DeepEqual(got, tc.expected) {
					t.Errorf("expected %v but got %v", tc.expected, got)
				}
			},
		)
	}
}

func TestCMap(t *testing.T) {
	ctx := context.Background()
	tests := []MapTestCase[int]{
		{
			"Square each element",
			[]int{1, 2, 3, 4, 5},
			func(n int) int { return n * n },
			[]int{1, 4, 9, 16, 25},
		},
		{
			"Return empty list when input is empty",
			[]int{},
			func(n int) int { return n * n },
			[]int{},
		},
	}

	for _, tc := range tests {
		t.Run(
			tc.name,
			func(t *testing.T) {
				got := make([]int, 0)
				supplierChan := make(chan int)
				go func() {
					defer close(supplierChan)
					for _, elem := range tc.A {
						supplierChan <- elem
					}
				}()
				for elem := range CMap(ctx, supplierChan, tc.mapFunc) {
					got = append(got, elem)
				}
				if !reflect.DeepEqual(got, tc.expected) {
					t.Errorf("expected %v but got %v", tc.expected, got)
				}
			},
		)
	}
}

func TestCMap_CancelCtx(t *testing.T) {
	GBufferSize = 0
	defer func() { GBufferSize = 64 }()
	ctx, cancel := context.WithCancel(context.Background())

	defer cancel()

	got := make([]int, 0)
	supplierChan, _ := Repeat(ctx, 1, -1)

	// Test will fail if cancel is ignored by CMap
	for elem := range CMap(ctx, supplierChan, func(n int) int { return n * n }) {
		cancel()
		_ = append(got, elem)
	}
}
