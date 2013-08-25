package goext

import (
	"testing"
)

func TestArraysAreEqualStrings_differentLengths(t *testing.T) {
	if ArraysAreEqualStrings([]string{}, []string{"Ruth"}) {
		t.Error("Error in comparing arrays, different lengths")
	}
}

func TestArraysAreEqualStrings_diffrentContents(t *testing.T) {
	if ArraysAreEqualStrings([]string{"Aaron"}, []string{"Ruth"}) {
		t.Error("Error in comparing arrays, different contents")
	}
}

func TestArraysAreEqualStrings_differentCase(t *testing.T) {
	if ArraysAreEqualStrings([]string{"ruth"}, []string{"Ruth"}) {
		t.Error("Error in comparing arrays, different case in strings")
	}
}

func TestArraysAreEqualStrings_bothEmpty(t *testing.T) {
	if !ArraysAreEqualStrings([]string{}, []string{}) {
		t.Error("Error in comparing arrays, both empty")
	}
}

func TestArraysAreEqualStrings_identical(t *testing.T) {
	if !ArraysAreEqualStrings([]string{"Ruth"}, []string{"Ruth"}) {
		t.Error("Error in comparing arrays, identical")
	}
}
