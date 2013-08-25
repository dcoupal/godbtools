package db

import (
	"reflect"
	"testing"
)

func TestCouchBaseProvider1(t *testing.T) {
	provider := GetDocProvider("couchbase://localhost/mybucket")
	providerType := reflect.TypeOf(provider).String()
	expectedType := "*db.CouchBase"
	if providerType != expectedType {
		t.Errorf("Invalid value for %s, expected %s, got %s", "Provider type", expectedType, providerType)
	}
}

func TestMongoDBProvider1(t *testing.T) {
	provider := GetDocProvider("mongodb://localhost/mydatabase/mycollection")
	providerType := reflect.TypeOf(provider).String()
	expectedType := "*db.MongoDB"
	if providerType != expectedType {
		t.Errorf("Invalid value for %s, expected %s, got %s", "Provider type", expectedType, providerType)
	}
	db := provider.get("DB")
	eDb := "mydatabase"
	if db != eDb {
		t.Errorf("Invalid value for %s, expected %s, got %s", "Database", eDb, db)
	}
	coll := provider.get("Coll")
	eColl := "mycollection"
	if coll != eColl {
		t.Errorf("Invalid value for %s, expected %s, got %s", "Collection", coll, eColl)
	}
}

func TestFileProvider1(t *testing.T) {
	provider := GetDocProvider("/path/to/my/file")
	providerType := reflect.TypeOf(provider).String()
	expectedType := "*db.TextFile"
	if providerType != expectedType {
		t.Errorf("Invalid value for %s, expected %s, got %s", "Provider type", expectedType, providerType)
	}
}
