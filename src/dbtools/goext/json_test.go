package goext

import (
	"testing"
)

func TestReadJsonFile1(t *testing.T) {
	doc := ReadJsonFile("test/data1.json")
	eDoc := map[string]interface{}{"name": "Guiness"}
	if !DocsAreEqual(doc, eDoc) {
		t.Errorf("Invalid value for %s, expected %v, got %v", "document", eDoc, doc)
	}
}

func TestReadJsonFile2(t *testing.T) {
	doc := ReadJsonFile("test/data2.json")
	eDoc := map[string]interface{}{"name": "Guiness", "country": "Ireland"}
	if !DocsAreEqual(doc, eDoc) {
		t.Errorf("Invalid value for %s, expected %v, got %v", "document", eDoc, doc)
	}
}

func TestReadJsonFile3(t *testing.T) {
	doc := ReadJsonFile("test/data3.json")
	required := []string{"ret"}
	properties := map[string]interface{}{
		"ret":  map[string]interface{}{"type": "number"},
		"date": map[string]interface{}{"type": "number"}}
	eDoc := map[string]interface{}{
		"$schema":    "http://json-schema.org/draft-03/schema",
		"id":         "http://jsonschema.net",
		"properties": properties,
		"required":   required}
	if !DocsAreEqual(doc, eDoc) {
		t.Errorf("Invalid value for %s, expected %v, got %v", "document", eDoc, doc)
	}
}
