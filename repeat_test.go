package gombinator

import (
	"context"
	"fmt"
	"testing"
)

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
