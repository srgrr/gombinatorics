package gombinator

import (
	"context"
	"fmt"
	"reflect"
	"testing"
)

type FilterTestCase[T any] struct {
	name      string
	A         []T
	criterion func(T) bool
	expected  []T
}

func TestFilter(t *testing.T) {
	ctx := context.Background()
	tests := []FilterTestCase[int]{
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

func TestCFilter(t *testing.T) {
	ctx := context.Background()
	tests := []FilterTestCase[int]{
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

func TestCFilter_CancelCtx(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	defer cancel()

	got := make([]int, 0)
	supplierChan := Range(ctx, 1, 1000000000)

	cancel()

	for elem := range CFilter(ctx, supplierChan, func(n int) bool { return n == 1000000000-1 }) {
		got = append(got, elem)
	}

	if len(got) != 0 {
		t.Errorf("Context cancellation either didn't happen or took too long")
	}
}

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
	supplierChan := Repeat(ctx, 1, -1)

	// Test will fail if cancel is ignored by CMap
	for elem := range CMap(ctx, supplierChan, func(n int) int { return n * n }) {
		cancel()
		_ = append(got, elem)
	}
}

func TestPartition(t *testing.T) {
	ctx := context.Background()
	strings := []string{"the", "quick", "brown", "fox", "jumps", "over", "the"}
	received := make([][]string, 0)

	partitionChannel, _ := Partition(ctx, strings, 2)

	for partition := range partitionChannel {
		received = append(received, partition)
	}
	expected := [][]string{
		{"the", "quick"},
		{"brown", "fox"},
		{"jumps", "over"},
		{"the"},
	}
	if !reflect.DeepEqual(expected, received) {
		t.Errorf("Expected %s\nReceived %s", expected, received)
	}
}

type PartitionWrongKValuesTestCase struct {
	name string
	k    int
}

func TestPartition_WrongKValues(t *testing.T) {
	ctx := context.Background()
	testCases := []PartitionWrongKValuesTestCase{
		{"k = 0", 0},
		{"k = -1", -1},
		{"k = -999", -999},
	}

	for _, testCase := range testCases {
		t.Run(
			testCase.name,
			func(t *testing.T) {
				strings := []string{"the", "quick", "brown", "fox", "jumps", "over", "the"}
				_, err := Partition(ctx, strings, testCase.k)
				if err == nil || err != ErrInvalidPartitionKValue {
					t.Errorf("Expected ErrInvalidPartitionKValue, received %+v instead", err)
				}
			},
		)
	}
}

func TestCPartition(t *testing.T) {
	ctx := context.Background()
	strings := []string{"the", "quick", "brown", "fox", "jumps", "over", "the"}
	received := make([][]string, 0)

	supplierChannel := make(chan string)
	go func() {
		defer close(supplierChannel)
		for _, str := range strings {
			supplierChannel <- str
		}
	}()

	partitionChannel, _ := Partition(ctx, strings, 2)

	for partition := range partitionChannel {
		received = append(received, partition)
	}

	expected := [][]string{
		{"the", "quick"},
		{"brown", "fox"},
		{"jumps", "over"},
		{"the"},
	}
	if !reflect.DeepEqual(expected, received) {
		t.Errorf("Expected %s\nReceived %s", expected, received)
	}
}

func TestCPartition_WrongKValues(t *testing.T) {
	ctx := context.Background()
	testCases := []PartitionWrongKValuesTestCase{
		{"k = 0", 0},
		{"k = -1", -1},
		{"k = -999", -999},
	}

	for _, testCase := range testCases {
		t.Run(
			testCase.name,
			func(t *testing.T) {
				ch := make(chan string)
				defer close(ch)
				_, err := CPartition(ctx, ch, testCase.k)
				if err == nil || err != ErrInvalidPartitionKValue {
					t.Errorf("Expected ErrInvalidPartitionKValue, received %+v instead", err)
				}
			},
		)
	}
}

func TestRange(t *testing.T) {
	ctx := context.Background()
	var total int
	for i := range Range(ctx, 1, 11) {
		total += i
	}
	if total != 10*11/2 {
		t.Errorf("Range didn't produce the right values")
	}
}

type RepeatTestCase struct {
	name string
	k    int
}

func TestERepeat_Finite(t *testing.T) {
	ctx := context.Background()
	testCases := []RepeatTestCase{
		{"k=1", 1},
		{"k=2", 2},
		{"k=100", 100},
	}
	for _, testCase := range testCases {
		t.Run(
			testCase.name,
			func(t *testing.T) {
				total := 0
				repeatChan, _ := ERepeat(ctx, 1, testCase.k)
				for range repeatChan {
					total++
				}
				if total != testCase.k {
					t.Errorf("Expected %d\nReceived %d\n", testCase.k, total)
				}
			},
		)
	}
}

func TestRepeat_Infinite(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	repeatChan, _ := ERepeat(ctx, 1, -1)
	cancel()
	total := 0
	for range repeatChan {
		total++
	}
	fmt.Println(total)
}

func TestERepeat_Error(t *testing.T) {
	ctx := context.Background()
	testCases := []RepeatTestCase{
		{"k=0", -3},
		{"k=1", -2},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name,
			func(t *testing.T) {
				_, err := ERepeat(ctx, 1, testCase.k)
				if err == nil {
					t.Errorf("Expected error")
				}
			},
		)
	}
}

func TestZip(t *testing.T) {
	ctx := context.Background()
	cities := []string{"london", "sf", "philly"}
	weather := []string{"cloudy", "foggy", "crazy"}

	zippedPairs := make([]Pair[string, string], 0)

	for pair := range Zip(ctx, cities, weather) {
		zippedPairs = append(zippedPairs, pair)
	}

	expected := []Pair[string, string]{
		{First: "london", Second: "cloudy"},
		{First: "sf", Second: "foggy"},
		{First: "philly", Second: "crazy"},
	}

	if !reflect.DeepEqual(expected, zippedPairs) {
		t.Errorf("Error:\nGot\t%+v\nExpected\t%+v", zippedPairs, expected)
	}
}

func TestCZip(t *testing.T) {
	ctx := context.Background()
	cities := Channelize(ctx, []string{"london", "sf", "philly"})
	weather := Channelize(ctx, []string{"cloudy", "foggy", "crazy", "remnant"})

	zippedPairs := make([]Pair[string, string], 0)

	for pair := range CZip(ctx, cities, weather) {
		zippedPairs = append(zippedPairs, pair)
	}

	expected := []Pair[string, string]{
		{First: "london", Second: "cloudy"},
		{First: "sf", Second: "foggy"},
		{First: "philly", Second: "crazy"},
	}

	if !reflect.DeepEqual(expected, zippedPairs) {
		t.Errorf("Error:\nGot\t%+v\nExpected\t%+v", zippedPairs, expected)
	}
}
