package functional

import (
	"context"
	"reflect"
	"testing"
)

func TestPartition(t *testing.T) {
	ctx := context.Background()
	strings := []string{"the", "quick", "brown", "fox", "jumps", "over", "the"}
	received := make([][]string, 0)
	for partition := range Partition(ctx, strings, 2) {
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

	for partition := range CPartition(ctx, supplierChannel, 2) {
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
