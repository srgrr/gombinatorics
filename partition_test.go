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
