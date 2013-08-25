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
	returnCode := doit([]string{"valjsondb", "-schema", "../../../testdata/schema_example1.json",
		"-data", "mongodb:/nosql2013/logs"})
	eReturnCode := 0
	if returnCode != eReturnCode {
		t.Errorf("Invalid value for %s, expected %s, got %s", "Return code", eReturnCode, returnCode)
	}

}

// Invalid documents, get a non-zero return code
func Test4(t *testing.T) {
	returnCode := doit([]string{"valjsondb", "-schema", "file:" + gopath + "/testdata/schema_example2.json",
		"-data", "mongodb:/nosql2013/logs", "-short"})
	eReturnCode := 1
	if returnCode != eReturnCode {
		t.Errorf("Invalid value for %s, expected %s, got %s", "Return code", eReturnCode, returnCode)
	}

}

func Test5(t *testing.T) {
	returnCode := doit([]string{"valjsondb", "-schema", "mongodb:/nosql2013/schemas?database=nosql2013&collection=logs",
		"-data", "mongodb:/nosql2013/logs", "-short"})
	eReturnCode := 1
	if returnCode != eReturnCode {
		t.Errorf("Invalid value for %s, expected %s, got %s", "Return code", eReturnCode, returnCode)
	}

}