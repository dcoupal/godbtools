package main

// Some system tests require a connection to MongoDB, others to CouchBase
// You need GOPATH to be set to the root of the workspace

import (
	"os"
	"testing"
)

var gopath string = os.Getenv("GOPATH")

//FIXME - the os.Exit on the help is exiting the test suite, failing it
func _Test1(t *testing.T) {
	returnCode := doit([]string{"valjsondb", "-help"})
	eReturnCode := 2
	if returnCode != eReturnCode {
		t.Errorf("Invalid value for %s, expected %s, got %s", "Return code", eReturnCode, returnCode)
	}

}

func Test2(t *testing.T) {
	returnCode := doit([]string{"valjsondb", "-version"})
	eReturnCode := 0
	if returnCode != eReturnCode {
		t.Errorf("Invalid value for %s, expected %s, got %s", "Return code", eReturnCode, returnCode)
	}

}

func Test3(t *testing.T) {
	returnCode := doit([]string{"valjsondb", "-schema", "file:" + gopath + "/testdata/schema_example3.json",
		"-data", "file:" + gopath + "/testdata/beers.data"})
	eReturnCode := 0
	if returnCode != eReturnCode {
		t.Errorf("Invalid value for %s, expected %s, got %s", "Return code", eReturnCode, returnCode)
	}

}

func Test4(t *testing.T) {
	returnCode := doit([]string{"valjsondb", "-schema", "couchbase:/buckets/test?1",
		"-data", "file:" + gopath + "/testdata/beers.data"})
	eReturnCode := 0
	if returnCode != eReturnCode {
		t.Errorf("Invalid value for %s, expected %s, got %s", "Return code", eReturnCode, returnCode)
	}

}

func Test5(t *testing.T) {
	returnCode := doit([]string{"valjsondb", "-schema", "couchbase:/buckets/test?3",
		"-data", "couchbase:/buckets/test/designs/dev_all/views/all"})
	eReturnCode := 0
	if returnCode != eReturnCode {
		t.Errorf("Invalid value for %s, expected %s, got %s", "Return code", eReturnCode, returnCode)
	}

}

func Test6(t *testing.T) {
	returnCode := doit([]string{"valjsondb", "-schema", "../../../testdata/schema_example1.json",
		"-data", "mongodb:/databases/nosql2013/collections/logs"})
	eReturnCode := 0
	if returnCode != eReturnCode {
		t.Errorf("Invalid value for %s, expected %s, got %s", "Return code", eReturnCode, returnCode)
	}

}

// Invalid documents, get a non-zero return code
func TestFailures1(t *testing.T) {
	returnCode := doit([]string{"valjsondb", "-schema", "file:" + gopath + "/testdata/schema_example2.json",
		"-data", "mongodb:/databases/nosql2013/collections/logs", "-short"})
	eReturnCode := 1
	if returnCode != eReturnCode {
		t.Errorf("Invalid value for %s, expected %s, got %s", "Return code", eReturnCode, returnCode)
	}

}

func TestFailures2(t *testing.T) {
	returnCode := doit([]string{"valjsondb", "-schema", "mongodb:/db/nosql2013/c/schemas?database=nosql2013&collection=logs",
		"-data", "mongodb:/databases/nosql2013/collections/logs", "-short"})
	eReturnCode := 1
	if returnCode != eReturnCode {
		t.Errorf("Invalid value for %s, expected %s, got %s", "Return code", eReturnCode, returnCode)
	}

}

func TestFailures3(t *testing.T) {
	returnCode := doit([]string{"valjsondb", "-schema", "couchbase:/buckets/test?2",
		"-data", "file:" + gopath + "/testdata/beers.data", "-short"})
	eReturnCode := 1
	if returnCode != eReturnCode {
		t.Errorf("Invalid value for %s, expected %s, got %s", "Return code", eReturnCode, returnCode)
	}

}
