package db

import (
	"fmt"
	"net/url"
)

import (
	"dbtools/goext"
)

type Doc map[string]interface{}

var supportedDbs = []string{"couchbase", "file", "mongodb"}

func IsDBSupported(dbname string) bool {
	return goext.StringInArray(dbname, supportedDbs)
}

/*
Valid URIs
  mongodb://localhost:2717/mydb/mycollection
  mongodb:/mydb/mycollection
  couchbase://localhost:27017/mybucket
  couchbase:/mybucket
  file:/path/to/my/file
  /path/to/my/file
*/
func GetDocProvider(uri string) (provider DocProvider) {
	provider = nil
	url, _ := url.Parse(uri)
	dbtype := url.Scheme
	host := url.Host
	path := url.Path
	query := url.RawQuery
	if dbtype == "couchbase" {
		provider = new(CouchBase)
		provider.SetDocProvider(host, path)
		provider.SetQuery(query)
	} else if dbtype == "mongodb" {
		provider = new(MongoDB)
		provider.SetDocProvider(host, path)
		provider.SetQuery(query)
	} else if dbtype == "file" || dbtype == "" {
		provider = new(TextFile)
		provider.SetDocProvider(host, path)
		provider.SetQuery(query)
	} else {
		fmt.Printf("No valid type in URI: %s", uri)
	}
	return
}

type DocProvider interface {
	Get(string) interface{}
	GetDoc(string) map[string]interface{}
	GetDocs() <-chan Doc
	GetQuery() string
	SetDocProvider(host string, path string)
	SetQuery(query string)
	//GetVersionedDoc() doc
}
