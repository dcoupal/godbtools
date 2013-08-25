// Package to extend Go functionality of maps
// The goal would be to not use this package on the long run, but rather
// rely on more standard packages available as 3rd party software

// TODO
//    - add the commas at the end of lines for the prettyMap to make it more JSON compatible.

package goext

import (
	"fmt"
	"sort"
)

func DocsAreEqual(aDoc map[string]interface{}, bDoc map[string]interface{}) bool {
	// TODO - add recursivity
	areEqual := true
	if len(aDoc) != len(bDoc) {
		areEqual = false
	} else {
		for k, _ := range aDoc {
			if value, ok := bDoc[k]; ok {
				switch value.(type) {
				case []interface{}:
					// TODO complete for array
				case string:
					if value != aDoc[k] {
						areEqual = false
					}
				}
			} else {
				areEqual = false
				break
			}
		}
	}
	return areEqual
}

func GetSortedKeys(aMap map[string]interface{}) []string {
	sKeys := make([]string, len(aMap))
	var i int = 0
	for k, _ := range aMap {
		sKeys[i] = k
		i++
	}
	sort.Strings(sKeys)
	return sKeys
}

func PrettyMap(data map[string]interface{}) string {
	return prettyMapInd(data, "")
}
func prettyMapInd(data map[string]interface{}, ind string) (res string) {
	// control indentation, by adding 2 spaces in recursive calls
	res = "{\n"
	nextInd := ind + "  "
	for _, key := range GetSortedKeys(data) {
		value := data[key]
		var line string
		switch value.(type) {
		case nil:
			line = key + ": null\n"
		case string:
			valueString := value.(string)
			line = key + `: "` + valueString + "\"\n"
		case int, int32, int64, float32, float64:
			line = fmt.Sprintf("%s: %v\n", key, value)
		case map[string]interface{}:
			valueMap := value.(map[string]interface{})
			line = key + ": " + prettyMapInd(valueMap, nextInd)
		case []interface{}:
			valueArray := value.([]interface{})
			line = key + ": " + prettyArrayInd(valueArray, nextInd)
		default:
			line = key + ": <unknown type>\n"
		}
		res = res + nextInd + line
	}
	res = res + ind + "}\n"
	return
}

func prettyArrayInd(data []interface{}, ind string) (res string) {
	res = "[\n"
	nextInd := ind + "  "
	for _, value := range data {
		var line string
		switch value.(type) {
		case nil:
			line = "null\n"
		case string:
			valueString := value.(string)
			line = "\"" + valueString + "\"\n"
		case int, int32, int64, float32, float64:
			line = fmt.Sprintf("%v\n", value)
		case map[string]interface{}:
			valueMap := value.(map[string]interface{})
			line = prettyMapInd(valueMap, nextInd)
		case []interface{}:
			valueArray := value.([]interface{})
			line = prettyArrayInd(valueArray, nextInd)
		default:
			line = "<unknown type>\n"
		}
		res = res + nextInd + line
	}
	res = res + ind + "]\n"
	return
}
