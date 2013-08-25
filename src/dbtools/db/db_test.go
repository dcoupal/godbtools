package db

import (
	"reflect"
	"testing"
)

func TestCouchBaseProvider1(t *testing.T) {
	if IsDBSupported("couchbase") {
		provider := GetDocProvider("couchbase://localhost/mybucket")
		providerType := reflect.TypeOf(provider).String()
		expectedType := "*db.CouchBase"
		if providerType != expectedType {
			t.Errorf("Invalid value for %s, expected %s, got %s", "Provider type", expectedType, providerType)
		}
		bucket := provider.Get("Bucket").(string)
		eBucket := "mybucket"
		if bucket != eBucket {
			t.Errorf("Invalid value for %s, expected %s, got %s", "Bucket", eBucket, bucket)
		}
	}
}

func TestMongoDBProvider1(t *testing.T) {
	if IsDBSupported("mongodb") {
		provider := GetDocProvider("mongodb://localhost/mydatabase/mycollection")
		providerType := reflect.TypeOf(provider).String()
		expectedType := "*db.MongoDB"
		if providerType != expectedType {
			t.Errorf("Invalid value for %s, expected %s, got %s", "Provider type", expectedType, providerType)
		}
		database := provider.Get("Database").(string)
		eDb := "mydatabase"
		if database != eDb {
			t.Errorf("Invalid value for %s, expected %s, got %s", "Database", eDb, database)
		}
		coll := provider.Get("Collection").(string)
		eColl := "mycollection"
		if coll != eColl {
			t.Errorf("Invalid value for %s, expected %s, got %s", "Collection", coll, eColl)
		}
	}
}

func TestFileProvider1(t *testing.T) {
	if IsDBSupported("mongodb") {
		provider := GetDocProvider("/path/to/my/file")
		providerType := reflect.TypeOf(provider).String()
		expectedType := "*db.TextFile"
		if providerType != expectedType {
			t.Errorf("Invalid value for %s, expected %s, got %s", "Provider type", expectedType, providerType)
		}
		path := provider.Get("Path").(string)
		ePath := "/path/to/my/file"
		if path != ePath {
			t.Errorf("Invalid value for %s, expected %s, got %s", "File", path, ePath)
		}
	}
}
