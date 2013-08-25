package goext

import ()

func ArraysAreEqualStrings(a1 []string, a2 []string) (areEqual bool) {
	areEqual = true
	if len(a1) != len(a2) {
		areEqual = false
	} else {
		for i, v := range a1 {
			if v != a2[i] {
				areEqual = false
				break
			}
		}
	}
	return
}
