package goext

import (
	"testing"
)

func TestGetSortedKeys1(t *testing.T) {
	aMap := map[string]interface{}{"McGuire": 70, "Maris": 61, "Ruth": 60}
	sKeys := GetSortedKeys(aMap)
	if !ArraysAreEqualStrings(sKeys, []string{"Maris", "McGuire", "Ruth"}) {
		t.Error("Error in sorting keys of map[string]int")
	}
}
