package gombinator

import (
	"context"
	"testing"
)

type RepeatTestCase struct {
	name string
	k    int
}

func TestRepeat_Error(t *testing.T) {
	ctx := context.Background()
	testCases := []RepeatTestCase{
		{"k=0", 0},
		{"k=1", -2},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name,
			func(t *testing.T) {
				_, err := Repeat(ctx, 1, testCase.k)
				if err == nil {
					t.Errorf("Expected error")
				}
			},
		)
	}
}
