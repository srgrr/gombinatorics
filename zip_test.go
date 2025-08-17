package gombinator

import (
	"context"
	"reflect"
	"testing"
)

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
