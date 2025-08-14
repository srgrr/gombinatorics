package functional

import (
	"context"
	"reflect"
	"testing"
)

func TestFilter(t *testing.T) {
	ctx := context.Background()

	numbers := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	evenNumbersChan := Filter(ctx, numbers, func(n int) bool { return n%2 == 0 })
	evenNumbers := make([]int, 0)
	for num := range evenNumbersChan {
		evenNumbers = append(evenNumbers, num)
	}
	expectedEven := []int{2, 4, 6, 8, 10}
	if !reflect.DeepEqual(evenNumbers, expectedEven) {
		t.Errorf("Expected %v, got %v", expectedEven, evenNumbers)
	}
}
