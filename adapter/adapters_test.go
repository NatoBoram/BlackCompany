package adapter_test

import (
	"slices"
	"testing"

	"github.com/NatoBoram/BlackCompany/adapter"
)

// TestToKeys_Empty checks if ToKeys returns an empty slice for an empty map.
func TestToKeys_Empty(t *testing.T) {
	empty := map[string]int{}
	got := adapter.ToKeys(empty)

	expected := []string{}
	if !slices.Equal(got, expected) {
		t.Errorf("ToKeys(empty) = %v, want %v", got, expected)
	}
}

// TestToKeys_Single checks if ToKeys returns the correct key for a single-entry
// map.
func TestToKeys_Single(t *testing.T) {
	single := map[string]int{"key": 42}
	got := adapter.ToKeys(single)

	expected := []string{"key"}
	if !slices.Equal(got, expected) {
		t.Errorf("ToKeys(single) = %v, want %v", got, expected)
	}
}

// TestToKeys_Multiple checks if ToKeys returns all keys for a multi-entry map.
func TestToKeys_Multiple(t *testing.T) {
	multiple := map[string]int{"first": 1, "second": 2, "third": 3}
	got := adapter.ToKeys(multiple)

	expected := []string{"first", "second", "third"}
	if !slices.Equal(got, expected) {
		t.Errorf("ToKeys(multiple) = %v, want %v", got, expected)
	}
}
