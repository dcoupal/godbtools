package db

import (
	"fmt"
	"testing"
)

import (
	"dbtools/goext"
	"labix.org/v2/mgo"
)

var verbose bool = false

func TestConnection(t *testing.T) {
	var host = "localhost:27017"
	session, err := mgo.Dial(host)
	if err != nil {
		t.Errorf("Can't connect with 'mgo' to %s\n", host)
	}
	defer session.Close()
}

func TestFindOne(t *testing.T) {
	var host = "localhost:27017"
	session, err := mgo.Dial(host)
	if err != nil {
		t.Errorf("Can't connect with 'mgo' to %s\n", host)
	}
	var doc map[string]interface{} = nil
	coll := session.DB("nosql2013").C("schemas")
	query := map[string]interface{}{"database": "nosql2013", "collection": "logs"}
	err = coll.Find(query).One(&doc) //debugger
	if verbose {
		fmt.Printf("%v", doc)
	}
	if _, ok := doc["_id"]; !ok {
		t.Errorf("Can't read one doc with 'mgo'")
	}
	defer session.Close()
}

func TestGetDoc(t *testing.T) {
	// TODO-FIXME if the collection name is wrong, you can't tell easily, need to wrap and diagnose errors like that.
	provider := GetDocProvider("mongodb://localhost:27017/nosql2013/schemas")
	doc := provider.GetDoc("database=nosql2013&collection=logs")
	vProperties := doc["properties"].(map[string]interface{})
	vReturn := vProperties["ret"].(map[string]interface{})
	value := vReturn["type"]
	eValue := "number"
	if value != eValue {
		t.Errorf("Invalid value for %s, expected %v, got %v", "Document value", eValue, value)
	}
}

func TestStr2mQuery(t *testing.T) {
	doc := Str2mQuery("database=nosql2013&collection=logs")
	eDoc := map[string]interface{}{"database": "nosql2013", "collection": "logs"}
	if !goext.DocsAreEqual(doc, eDoc) {
		t.Errorf("Invalid value for %s, expected %v, got %v", "document", eDoc, doc)
	}
}
