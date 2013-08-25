package goext

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

func ReadJsonFile(filename string) map[string]interface{} {
	var doc map[string]interface{} = nil
	fileContents, e := ioutil.ReadFile(filename)
	if e != nil {
		fmt.Printf("Error reading file: %v\n", e)
		os.Exit(1)
	}
	json.Unmarshal([]byte(fileContents), &doc)
	return doc
}
