package set

import (
	"testing"
)

func Test_union(t *testing.T) {
	s1 := Set{"hello": true, "world": true}
	s2 := Set{"hello": true, "Nathan": true}

	actual := s1.Union(s2)
	expectedUnion := Set{"hello": true, "world": true, "Nathan": true}
	if actual.Len() != expectedUnion.Len() {
		t.Fatal("Expected", expectedUnion.Len(), "but was", actual.Len())
	}

	if !actual.Contains("hello") {
		t.Fatal("Set does not contain expected element.")
	}

	if !expectedUnion.Equal(actual) {
		t.Fatal("Sets are not equal.")
	}
}

func Test_diff(t *testing.T) {
	s1 := Set{"hello": true, "world": true}
	s2 := Set{"hello": true, "Nathan": true}

	actual := s1.Diff(s2)
	expectedUnion := Set{"world": true}
	if len(actual) != len(expectedUnion) {
		t.Fatal("Expected", len(expectedUnion), "but was", len(actual))
	}

	for m := range expectedUnion {
		if !actual[m] {
			t.Fatal(m, " not found in set.")
		}
	}
}

func Test_intersection(t *testing.T) {
	s1 := Set{"hello": true, "world": true}
	s2 := Set{"hello": true, "Nathan": true}

	actual := s1.Intersection(s2)
	expectedUnion := Set{"hello": true}
	if len(actual) != len(expectedUnion) {
		t.Fatal("Expected", len(expectedUnion), "but was", len(actual))
	}

	for m := range expectedUnion {
		if !actual[m] {
			t.Fatal(m, " not found in set.")
		}
	}
}
