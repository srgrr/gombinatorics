package gombinator

import (
	"context"
	"testing"
)

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
