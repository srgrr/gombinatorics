package gombinator

import (
	"context"
	"fmt"
	"reflect"
	"testing"
)

var BENCHMARK_LIMIT = 10000000
var BENCHMARK_MOD int = 1e9 + 7
var BENCHMARK_TARGET int = 3691500

type FilterTestCase[T any] struct {
	name      string
	A         []T
	criterion func(T) bool
	expected  []T
}

type MapTestCase[T any] struct {
	name     string
	A        []T
	mapFunc  func(T) T
	expected []T
}

type PartitionWrongKValuesTestCase struct {
	name string
	k    int
}

type RepeatTestCase struct {
	name string
	k    int
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

// Benchmarks

func BenchmarkRangeMapFilter_Gombinator(b *testing.B) {
	b.ReportAllocs()
	b.SetBytes(1)
	for n := 0; n < b.N; n++ {
		var total int = 0
		ctx := context.Background()
		squaredEvenNumbers :=
			CMap(ctx,
				CFilter(ctx,
					Range(ctx, 0, BENCHMARK_LIMIT),
					func(n int) bool { return n%2 == 0 },
				),
				func(n int) int { return n * n },
			)
		// Consume the channel for benchmarking
		for elem := range squaredEvenNumbers {
			total = (total + elem) % BENCHMARK_MOD
		}
		if total != BENCHMARK_TARGET {
			b.Errorf("Benchmark doesn't compute the right value")
		}
	}
}

func getRangeChan(l, r int) <-chan int {
	ch := make(chan int, 64)
	go func() {
		defer close(ch)
		for i := l; i < r; i++ {
			ch <- i
		}
	}()
	return ch
}

func getMapChan(inputChan <-chan int, f func(n int) int) <-chan int {
	ch := make(chan int, 64)
	go func() {
		defer close(ch)
		for elem := range inputChan {
			ch <- f(elem)
		}
	}()
	return ch
}

func getFilterChan(inputChan <-chan int, f func(n int) bool) <-chan int {
	ch := make(chan int, 64)
	go func() {
		defer close(ch)
		for elem := range inputChan {
			if f(elem) {
				ch <- elem
			}
		}
	}()
	return ch
}

func BenchmarkRangeMapFilter_Vanilla(b *testing.B) {
	b.ReportAllocs()
	b.SetBytes(1)
	for n := 0; n < b.N; n++ {
		var total int = 0
		squaredEvenNumbers :=
			getMapChan(
				getFilterChan(
					getRangeChan(0, BENCHMARK_LIMIT),
					func(n int) bool { return n%2 == 0 },
				),
				func(n int) int { return n * n },
			)
		// Consume the channel for benchmarking
		for elem := range squaredEvenNumbers {
			total = (total + elem) % BENCHMARK_MOD
		}
		if total != BENCHMARK_TARGET {
			b.Errorf("Benchmark doesn't compute the right value")
		}
	}
}

func getRangeChanCtx(ctx context.Context, l, r int) <-chan int {
	ch := make(chan int, 64)
	go func() {
		defer close(ch)
		for i := l; i < r; i++ {
			select {
			case <-ctx.Done():
				return
			default:
				ch <- i
			}
		}
	}()
	return ch
}

func getMapChanCtx(ctx context.Context, inputChan <-chan int, f func(n int) int) <-chan int {
	ch := make(chan int, 64)
	go func() {
		defer close(ch)
		for elem := range inputChan {
			select {
			case <-ctx.Done():
				return
			default:
				ch <- f(elem)
			}
		}
	}()
	return ch
}

func getFilterChanCtx(ctx context.Context, inputChan <-chan int, f func(n int) bool) <-chan int {
	ch := make(chan int, 64)
	go func() {
		defer close(ch)
		for elem := range inputChan {
			select {
			case <-ctx.Done():
				return
			default:
				if f(elem) {
					ch <- elem
				}
			}
		}
	}()
	return ch
}

func BenchmarkRangeMapFilter_VanillaContext(b *testing.B) {
	b.ReportAllocs()
	b.SetBytes(1)
	ctx := context.Background()
	for n := 0; n < b.N; n++ {
		var total int = 0
		squaredEvenNumbers :=
			getMapChanCtx(ctx,
				getFilterChanCtx(ctx,
					getRangeChanCtx(ctx, 0, BENCHMARK_LIMIT),
					func(n int) bool { return n%2 == 0 },
				),
				func(n int) int { return n * n },
			)
		// Consume the channel for benchmarking
		for elem := range squaredEvenNumbers {
			total = (total + elem) % BENCHMARK_MOD
		}
		if total != BENCHMARK_TARGET {
			b.Errorf("Benchmark doesn't compute the right value")
		}
	}
}

func Benchmark_Sequential(b *testing.B) {
	b.ReportAllocs()
	b.SetBytes(1)
	for n := 0; n < b.N; n++ {
		total := 0
		for i := 0; i < BENCHMARK_LIMIT; i++ {
			if i%2 == 0 {
				total = (total + i*i) % BENCHMARK_MOD
			}
		}
		if total != BENCHMARK_TARGET {
			b.Errorf("Benchmark doesn't compute the right value")
		}
	}
}
