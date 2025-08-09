package gombinatorics

import (
	"reflect"
	"testing"
)

type TestCase[P any, Q any] struct {
	name     string
	A        []P
	B        []Q
	expected []Pair[P, Q]
}

func TestCartesianProduct(t *testing.T) {
	tests := []TestCase[string, string]{
		{
			"Raccoons and rats vs cheese and trash",
			[]string{"🦝", "🐀"},
			[]string{"🧀", "🗑️"},
			[]Pair[string, string]{
				{"🦝", "🧀"},
				{"🦝", "🗑️"},
				{"🐀", "🧀"},
				{"🐀", "🗑️"},
			},
		},
		{
			"Empty list left",
			[]string{},
			[]string{"🧀", "🗑️"},
			[]Pair[string, string]{},
		},
		{
			"Empty list right",
			[]string{"🦝", "🐀"},
			[]string{},
			[]Pair[string, string]{},
		},
		{
			"Empty lists",
			[]string{},
			[]string{},
			[]Pair[string, string]{},
		},
	}

	for _, tc := range tests {
		t.Run(
			tc.name,
			func(t *testing.T) {
				prodLen := len(tc.A) * len(tc.B)
				got := make([]Pair[string, string], 0, prodLen)
				for obtainedPair := range CartesianProduct(tc.A, tc.B) {
					got = append(got, obtainedPair)
				}
				if !reflect.DeepEqual(got, tc.expected) {
					t.Errorf("Error:\nGot\t%+v\nExpected\t%+v", got, tc.expected)
				}
			},
		)
	}
}

func TestSubsets(t *testing.T) {
	nums := []int{1, 2, 4, 8}
	total := 0
	numSubsets := 0
	for subset := range Subsets(nums) {
		numSubsets++
		for _, elem := range subset {
			total += elem
		}
	}
	if total != 8*(1+2+4+8) || numSubsets != 16 {
		t.Errorf("%d %d", total, numSubsets)
	}
}

func TestSubsetsOfFixedSize(t *testing.T) {
	nums := []int{1, 2, 4, 8}
	total := 0
	numSubsets := 0
	for subset := range SubsetsOfFixedSize(nums, 2) {
		numSubsets++
		for _, elem := range subset {
			total += elem
		}
	}
	if total != 3*(1+2+4+8) || numSubsets != 6 {
		t.Errorf("%d %d", total, numSubsets)
	}
}

func TestZip(t *testing.T) {
	cities := []string{"london", "sf", "philly"}
	weather := []string{"cloudy", "foggy", "crazy"}

	zippedPairs := make([]Pair[string, string], 0)

	for pair := range Zip(cities, weather) {
		zippedPairs = append(zippedPairs, pair)
	}

	expected := []Pair[string, string]{
		{"london", "cloudy"},
		{"sf", "foggy"},
		{"philly", "crazy"},
	}

	if !reflect.DeepEqual(expected, zippedPairs) {
		t.Errorf("Error:\nGot\t%+v\nExpected\t%+v", zippedPairs, expected)
	}
}
