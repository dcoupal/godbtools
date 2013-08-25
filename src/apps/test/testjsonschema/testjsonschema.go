package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

import (
	"github.com/sigu-399/gojsonschema"
)

func get_json(format string) interface{} {
	var jsonToValidate interface{}
	var err error
	if format == "file" {
		// (gdb) print jsonToValidate
		// $1 = {_type = 0x1c5000 <egcdata+186880>, data = 0xc2000a89c0}
		jsonToValidate, err = GetFileJson("/Users/dcoupal/json/data_example1.json")
	} else if format == "json" {
		jsonToValidate = "{\"student_id\":33717,\"type\":\"homework\",\"score\":20}"
	} else if format == "map" {
		doc := map[string]interface{}{
			"student_id": 33717,
			"type":       "homework",
			"score":      20,
		}
		//var jsonDoc interface{}
		//jsonToValidate, err = json.Marshal(doc)
		jsonToValidate = doc
		//jsonToValidate = interface{}(jsonDoc)
		if err != nil {
			panic(err.Error())
		}
	}
	return (jsonToValidate)
}

func main() {
	schema, err := gojsonschema.NewJsonSchemaDocument("file:///Users/dcoupal/json/schema_example1.json")

	if err != nil {
		panic(err.Error())
	}

	jsonToValidate := get_json("map")
	validationResult := schema.Validate(jsonToValidate)
	fmt.Printf("IsValid %v\n", validationResult.IsValid())
	fmt.Printf("%v\n", validationResult.GetErrorMessages())

	os.Exit(0)
}

// From gojsonschema
// Helper function to read a json from a http request
func GetFileJson(filepath string) (interface{}, error) {

	bodyBuff, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	var document interface{}
	err = json.Unmarshal(bodyBuff, &document)
	if err != nil {
		return nil, err
	}

	return document, nil
}
