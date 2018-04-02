package main

import "testing"

func TestFilterKeeps(t *testing.T) {
	vv := []string{
		"a",
		"b",
	}
	f := func(s string) bool {
		return true
	}

	actual := Filter(vv, f)
	if actual[0] != "a" || actual[1] != "b" || len(actual) != 2 {
		t.Fatalf("filter didnt keep: %v", actual)
	}
}

// TestFilterDiscards
func TestFilterDiscards(t *testing.T) {
	vv := []string{
		"a",
		"b",
	}
	f := func(s string) bool {
		if s == "a" {
			return true
		}
		return false
	}

	actual := Filter(vv, f)
	if actual[0] != "a" || len(actual) != 1 {
		t.Fatalf("filter didnt discard: %v", actual)
	}
}
