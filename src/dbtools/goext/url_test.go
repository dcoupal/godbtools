package goext

import (
	"net/url"
	"testing"
)

func TestCouchBaseURI1(t *testing.T) {
	var expect string
	u, _ := url.Parse("couchbase://localhost/mybucket")
	expect = "couchbase"
	if u.Scheme != expect {
		t.Errorf("Invalid value for %s, expected %s, got %s", "scheme", expect, u.Scheme)
	}
	expect = "localhost"
	if u.Host != expect {
		t.Errorf("Invalid value for %s, expected %s, got %s", "host", expect, u.Host)
	}
	expect = "/mybucket"
	if u.Path != expect {
		t.Errorf("Invalid value for %s, expected %s, got %v", "bucket", expect, u.Path)
	}
}

func TestCouchBaseURI2(t *testing.T) {
	var expect string
	u, _ := url.Parse("couchbase:/mybucket")
	expect = "couchbase"
	if u.Scheme != expect {
		t.Errorf("Invalid value for %s, expected %s, got %s", "scheme", expect, u.Scheme)
	}
	expect = "/mybucket"
	if u.Path != expect {
		t.Errorf("Invalid value for %s, expected %s, got %v", "bucket", expect, u.Path)
	}
}

func TestMongoURI1(t *testing.T) {
	var expect string
	u, _ := url.Parse("mongodb://localhost:27017/mydb/mycollection")
	expect = "mongodb"
	if u.Scheme != expect {
		t.Errorf("Invalid value for %s, expected %s, got %s", "scheme", expect, u.Scheme)
	}
	expect = "localhost:27017"
	if u.Host != expect {
		t.Errorf("Invalid value for %s, expected %s, got %s", "host", expect, u.Host)
	}
	expect = "/mydb/mycollection"
	if u.Path != expect {
		t.Errorf("Invalid value for %s, expected %s, got %v", "bucket", expect, u.Path)
	}
}

func TestMongoURI2(t *testing.T) {
	var expect string
	u, _ := url.Parse("mongodb:/mydb/mycollection")
	expect = "mongodb"
	if u.Scheme != expect {
		t.Errorf("Invalid value for %s, expected %s, got %s", "scheme", expect, u.Scheme)
	}
	expect = "/mydb/mycollection"
	if u.Path != expect {
		t.Errorf("Invalid value for %s, expected %s, got %v", "bucket", expect, u.Path)
	}
}

func TestMongoURI3(t *testing.T) {
	var expect string
	u, _ := url.Parse("mongodb:/mydb/mycollection?db=mydb&coll=mycoll")
	expect = "mongodb"
	if u.Scheme != expect {
		t.Errorf("Invalid value for %s, expected %s, got %s", "scheme", expect, u.Scheme)
	}
	expect = "/mydb/mycollection"
	if u.Path != expect {
		t.Errorf("Invalid value for %s, expected %s, got %v", "bucket", expect, u.Path)
	}
	expect = "db=mydb&coll=mycoll"
	if u.RawQuery != expect {
		t.Errorf("Invalid value for %s, expected %s, got %v", "query", expect, u.RawQuery)
	}
	q, _ := url.ParseQuery(u.RawQuery)
	dbs := []string{"mydb"}
	colls := []string{"mycoll"}
	expectQ := url.Values{"db": dbs, "coll": colls}
	if !ValuesAreEqual(q, expectQ) {
		t.Errorf("Invalid value for %s, expected %v, got %v", "values in query", expectQ, q)
	}
}

func TestFileURI1(t *testing.T) {
	var expect string
	u, _ := url.Parse("file:/path/to/file")
	expect = "file"
	if u.Scheme != expect {
		t.Errorf("Invalid value for %s, expected %s, got %s", "scheme", expect, u.Scheme)
	}
	expect = "/path/to/file"
	if u.Path != expect {
		t.Errorf("Invalid value for %s, expected %s, got %v", "bucket", expect, u.Path)
	}
}

func TestFileURI2(t *testing.T) {
	var expect string
	u, _ := url.Parse("/path/to/file")
	expect = ""
	if u.Scheme != expect {
		t.Errorf("Invalid value for %s, expected %s, got %s", "scheme", expect, u.Scheme)
	}
	expect = "/path/to/file"
	if u.Path != expect {
		t.Errorf("Invalid value for %s, expected %s, got %v", "bucket", expect, u.Path)
	}
}

func TestFileURI3(t *testing.T) {
	var expect string
	u, _ := url.Parse("../path/to/file")
	expect = ""
	if u.Scheme != expect {
		t.Errorf("Invalid value for %s, expected %s, got %s", "scheme", expect, u.Scheme)
	}
	expect = "../path/to/file"
	if u.Path != expect {
		t.Errorf("Invalid value for %s, expected %s, got %v", "bucket", expect, u.Path)
	}
}

func FailingTestFileURI4(t *testing.T) {
	var expect string
	u, _ := url.Parse("file:../path/to/file")
	expect = "file"
	if u.Scheme != expect {
		t.Errorf("Invalid value for %s, expected %s, got %s", "scheme", expect, u.Scheme)
	}
	expect = "../path/to/file"
	if u.Path != expect {
		t.Errorf("Invalid value for %s, expected %s, got %v", "bucket", expect, u.Path)
	}
}

func TestGetParts1(t *testing.T) {
	var uri, match, expected string
	var matches []string
	var ok bool
	var sections [][]string
	uri = "/_design/mydesign/_view/myview"
	sections = [][]string{{""}, {"_design", "designs"}, {}, {"_view", "views"}, {}}
	if _, ok = GetParts(uri, sections); !ok {
		t.Errorf("Should match %v, %v", uri, sections)
	}
	uri = "/designs/mydesign/views/myview"
	if matches, ok = GetParts(uri, sections); !ok {
		t.Errorf("Should match %v, %v", uri, sections)
	}
	match = matches[2]
	expected = "mydesign"
	if match != expected {
		t.Errorf("Should match %v, %v", match, expected)
	}
	match = matches[4]
	expected = "myview"
	if match != expected {
		t.Errorf("Should match %v, %v", match, expected)
	}
	uri = "/designs/mydesign/view/myview"
	if _, ok = GetParts(uri, sections); ok {
		t.Errorf("Should not match %v, %v", uri, sections)
	}
}

func TestGetParts2(t *testing.T) {
	var uri, match, expected string
	var matches []string
	var ok bool
	var sections [][]string
	uri = "/buckets/mybucket"
	sections = [][]string{{""}, {"buckets"}, {}}
	if matches, ok = GetParts(uri, sections); !ok {
		t.Errorf("Should match %v, %v", uri, sections)
	}
	match = matches[2]
	expected = "mybucket"
	if match != expected {
		t.Errorf("Should match %v, %v", match, expected)
	}
}
