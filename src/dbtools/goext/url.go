package goext

import (
	"net/url"
	"strings"
)

func GetParts(path string, sections [][]string) ([]string, bool) {
	matches := make([]string, len(sections))
	var ok = true
	pathParts := strings.Split(path, "/")
	if len(pathParts) == len(sections) {
		for i, section := range sections {
			if len(section) == 0 {
				matches[i] = pathParts[i]
			} else {
				submatch := false
				for _, subsection := range section {
					if subsection == pathParts[i] {
						submatch = true
						break
					}
				}
				if submatch == false {
					ok = false
					break
				}
			}
		}
	} else {
		ok = false
	}
	if ok == false {
		matches = []string{}
	}
	return matches, ok
}

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
