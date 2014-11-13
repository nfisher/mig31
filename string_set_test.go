package main

import (
	"testing"
)

func Test_inequality(t *testing.T) {
	s1 := StringSet{"hello": true, "world": true}
	s2 := StringSet{"hello": true, "Nathan": true}

	if s1.Equal(s2) {
		t.Fatal("Sets should not be equal but were", s1)
	}
}

func Test_union(t *testing.T) {
	s1 := StringSet{"hello": true, "world": true}
	s2 := StringSet{"hello": true, "Nathan": true}

	actual := s1.Union(s2)
	expectedUnion := StringSet{"hello": true, "world": true, "Nathan": true}
	if actual.Len() != expectedUnion.Len() {
		t.Fatal("Expected", expectedUnion.Len(), "but was", actual.Len())
	}

	if !actual.Contains("hello") {
		t.Fatal("StringSet does not contain expected element.")
	}

	if !expectedUnion.Equal(actual) {
		t.Fatal("Sets are not equal.")
	}
}

func Test_diff(t *testing.T) {
	s1 := StringSet{"hello": true, "world": true}
	s2 := StringSet{"hello": true, "Nathan": true}

	actual := s1.Diff(s2)
	expectedUnion := StringSet{"world": true}
	if len(actual) != len(expectedUnion) {
		t.Fatal("Expected", len(expectedUnion), "but was", len(actual))
	}

	for m := range expectedUnion {
		if !actual[m] {
			t.Fatal(m, " not found in StringSet.")
		}
	}
}

func Test_intersection(t *testing.T) {
	s1 := StringSet{"hello": true, "world": true}
	s2 := StringSet{"hello": true, "Nathan": true}

	actual := s1.Intersection(s2)
	expectedUnion := StringSet{"hello": true}
	if len(actual) != len(expectedUnion) {
		t.Fatal("Expected", len(expectedUnion), "but was", len(actual))
	}

	for m := range expectedUnion {
		if !actual[m] {
			t.Fatal(m, " not found in StringSet.")
		}
	}
}
