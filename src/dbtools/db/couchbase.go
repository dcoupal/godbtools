package db

import (
	"reflect"
)
type CouchBase struct {
	dbtype string
	host   string
	port   int
	bucket string
	query  string
}

func (o *CouchBase) get(k string) interface{} {
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
	//
}

func (o *CouchBase) SetQuery(query string) {
	o.query = query
}
