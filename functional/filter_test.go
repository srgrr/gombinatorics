package functional

import (
	"context"
	"reflect"
	"testing"
)

type TestCase[T any] struct {
	name      string
	A         []T
	criterion func(T) bool
	expected  []T
}

func TestFilter(t *testing.T) {
	ctx := context.Background()
	tests := []TestCase[int]{
		{
			"Filter even numbers",
			[]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			func(n int) bool { return n%2 == 0 },
			[]int{2, 4, 6, 8, 10},
		},
		{
			"Return empy list when no elements match",
			[]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			func(n int) bool { return n == 0 },
			[]int{},
		},
	}

	for _, tc := range tests {
		t.Run(
			tc.name,
			func(t *testing.T) {
				got := make([]int, 0)
				for elem := range Filter(ctx, tc.A, tc.criterion) {
					got = append(got, elem)
				}
				if !reflect.DeepEqual(got, tc.expected) {
					t.Errorf("expected %v but got %v", tc.expected, got)
				}
			},
		)
	}
}

func TestFilter_DoneCtx(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel the context immediately

	got := make([]int, 0)
	for elem := range Filter(ctx, []int{1, 2, 3}, func(n int) bool { return n%2 == 0 }) {
		got = append(got, elem)
	}

	if len(got) != 0 {
		t.Errorf("expected no elements but got %v", got)
	}
}

func TestCFilter(t *testing.T) {
	ctx := context.Background()
	tests := []TestCase[int]{
		{
			"Filter even numbers",
			[]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			func(n int) bool { return n%2 == 0 },
			[]int{2, 4, 6, 8, 10},
		},
		{
			"Return empy list when no elements match",
			[]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			func(n int) bool { return n == 0 },
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
				for elem := range CFilter(ctx, supplierChan, tc.criterion) {
					got = append(got, elem)
				}
				if !reflect.DeepEqual(got, tc.expected) {
					t.Errorf("expected %v but got %v", tc.expected, got)
				}
			},
		)
	}
}

func TestCFilter_DoneCtx(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel the context immediately

	got := make([]int, 0)
	supplierChan := make(chan int)
	go func() {
		defer close(supplierChan)
		for _, elem := range []int{1, 2, 3} {
			supplierChan <- elem
		}
	}()

	for elem := range CFilter(ctx, supplierChan, func(n int) bool { return n%2 == 0 }) {
		got = append(got, elem)
	}

	if len(got) != 0 {
		t.Errorf("expected no elements but got %v", got)
	}
}
