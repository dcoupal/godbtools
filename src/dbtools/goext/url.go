package goext

import (
	"net/url"
)

func ValuesAreEqual(aDoc url.Values, bDoc url.Values) bool {
	// TODO - add recursivity
	areEqual := true
	if len(aDoc) != len(bDoc) {
		areEqual = false
	} else {
		for k, aValue := range aDoc {
			if bValue, ok := bDoc[k]; ok {
				if len(aValue) == len(bValue) {
					// compare the 2 arrays
					for i, value := range aValue {
						if value != bValue[i] {
							areEqual = false
							break
						}
					}
				} else {
					areEqual = false
					break
				}
			} else {
				areEqual = false
				break
			}
		}
	}
	return areEqual
}
