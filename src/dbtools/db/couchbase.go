package db

import (
	"reflect"
	"strings"
)

type CouchBase struct {
	dbtype string
	Host   string
	Bucket string
	port   int
	query  string
}

func (o *CouchBase) Get(k string) interface{} {
	var v interface{} = reflect.ValueOf(o).Elem().FieldByName(k).Interface()
	return v
}

func (o *CouchBase) GetDoc(query string) (doc map[string]interface{}) {
	return
}

func (o *CouchBase) GetQuery() string {
	return o.query
}

func (o *CouchBase) SetDocProvider(host string, path string) {
	o.Host = host
	arr := strings.Split(path, "/") // First item is empty
	if len(arr) != 2 {
		panic("The path to the database should have 1 field, the bucket name")
	}
	o.Bucket = arr[1]
}

func (o *CouchBase) SetQuery(query string) {
	o.query = query
}
