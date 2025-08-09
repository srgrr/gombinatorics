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
