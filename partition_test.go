package gombinator

import (
	"context"
	"reflect"
	"testing"
)

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
