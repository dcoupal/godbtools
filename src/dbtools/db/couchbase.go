package db

import (
	"fmt"
	"reflect"
	"strings"
)

import (
	"github.com/couchbaselabs/go-couchbase"
)

type CouchBase struct {
	dbtype string
	Host   string
	Bucket string
	Pool   string
	port   int
	Query  string
}

func (o *CouchBase) Get(k string) interface{} {
	var v interface{} = reflect.ValueOf(o).Elem().FieldByName(k).Interface()
	return v
}

func (o *CouchBase) GetDoc(query string) map[string]interface{} {
	var doc map[string]interface{} = nil

	c, err := couchbase.Connect("http://" + o.Host)
	mf(err, "connect - "+o.Host)

	p, err := c.GetPool(o.Pool)
	mf(err, "pool")

	b, err := p.GetBucket(o.Bucket)
	mf(err, "bucket")
	defer b.Close()

	err = b.Get(o.Query, &doc)
	mf(err, "bucket.get")

	return doc
}

func (o *CouchBase) GetDocs() <-chan Doc {
	ch := make(chan Doc)
	return ch
}

func (o *CouchBase) GetQuery() string {
	return o.Query
}

func (o *CouchBase) SetDocProvider(host string, path string) {
	if o.Host != "" {
		o.Host = host
	} else {
		o.Host = "localhost:8091"
	}
	arr := strings.Split(path, "/") // First item is empty
	if len(arr) == 2 {
		o.Pool = "default"
		o.Bucket = arr[1]
	} else if len(arr) == 3 {
		o.Pool = arr[1]
		o.Bucket = arr[2]
	} else {
		panic("The path to the database should have 1 field, the bucket name")
	}
}

func (o *CouchBase) SetQuery(query string) {
	o.Query = query
}

func mf(err error, msg string) {
	if err != nil {
		//log.Fatalf("%v: %v", msg, err)
		fmt.Printf("%v", msg)
		panic(err)
	}
}
